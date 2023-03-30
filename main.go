package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gummy789j/telegram-quote-bot/internal/command"
	"github.com/gummy789j/telegram-quote-bot/internal/constant"
	"github.com/gummy789j/telegram-quote-bot/internal/domain"
	"github.com/shopspring/decimal"
)

var defaultCli *http.Client
var token string
var latestUpdateID int64

var (
	celebrationEmoji = "&#127882;"
)

type exchange string

var (
	MAX   exchange = "MAX"
	Rybit exchange = "Rybit"
)

func init() {
	defaultCli = http.DefaultClient
	defaultCli.Timeout = 30 * time.Second

	token = os.Getenv("TELEGRAM_BOT_TOKEN")

	if len(token) == 0 {
		panic("token is empty")
	}

	// get the latest update id and store it
	umResp, err := getUpdateMessage()
	if err != nil {
		panic(err)
	}

	if len(umResp.Result) == 0 {
		return
	}

	for _, v := range umResp.Result {
		if v.UpdateID > latestUpdateID {
			latestUpdateID = v.UpdateID
		}
	}
}

var wg = &sync.WaitGroup{}

func main() {
	var err error

	defer func() {
		errorNotify(err.Error())
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 365*24*time.Hour)
	defer cancel()

	wg.Add(2)

	// notify arbitrage job
	go notifyArbitrage(ctx)

	// reply commands job
	go replyCommand(ctx)

	fmt.Println("jobs running")
	wg.Wait()

}

func replyCommand(ctx context.Context) {
	ticker := time.NewTicker(constant.ReplyFreq)
	defer ticker.Stop()
	fmt.Println("reply job start")
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println(strings.Repeat("=", 20))
			fmt.Println("reply command work done")
			fmt.Println(strings.Repeat("=", 20))
			return
		case <-ticker.C:
		}

		umResp, err := getUpdateMessage()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if len(umResp.Result) == 0 {
			continue
		}

		var newLastUpdateID = latestUpdateID
		for _, v := range umResp.Result {
			if v.UpdateID <= latestUpdateID {
				continue
			}

			if v.UpdateID > newLastUpdateID {
				newLastUpdateID = v.UpdateID
			}

			if v.Message == nil {
				continue
			}

			if v.Message.Text == nil {
				continue
			}

			if len(v.Message.Entities) == 0 {
				continue
			}

			if v.Message.Entities[0].Type != "bot_command" {
				continue
			}

			cHandler := commandFactory(constant.CommandType(strings.TrimPrefix(*v.Message.Text, "/")))
			if cHandler == nil {
				continue
			}

			if v.Message.From == nil {
				continue
			}

			replyMsg := cHandler.Reply(v.Message.From.ID)

			if len(replyMsg) == 0 {
				continue
			}

			if err := botSendMessage(botSendMessageReq{
				chatID: v.Message.Chat.ID,
				msg:    replyMsg,
			}); err != nil {
				log.Println(err.Error())
			}
		}
		latestUpdateID = newLastUpdateID
	}

}

func notifyArbitrage(ctx context.Context) {

	ticker := time.NewTicker(constant.NotifyFreq)
	fmt.Println("notify job start")

	defer ticker.Stop()

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println(strings.Repeat("=", 20))
			fmt.Println("notify arbitrage work done")
			fmt.Println(strings.Repeat("=", 20))
			return
		case <-ticker.C:
		}
		qInfo, err := fetchExchangeQuotation([]exchange{MAX, Rybit})
		if err != nil {
			log.Println(err.Error())
			return
		}

		aInfo := calArbitrageInfo(constant.DefaultInvest, qInfo[Rybit].buyPrice, qInfo[MAX].sellPrice)

		if aInfo.Arbitrage.LessThan(constant.MinArbitrage) {
			return
		}

		if aInfo.Spread.LessThan(constant.MinSpread) {
			return
		}

		var isExcitedArbitrage, isExcitedSpread bool

		if aInfo.Arbitrage.GreaterThanOrEqual(constant.ExcitedArbitrage) {
			isExcitedArbitrage = true
		}

		if aInfo.Spread.GreaterThanOrEqual(constant.ExcitedSpread) {
			isExcitedSpread = true
		}

		if err = botSendNotifyMessage(botSendNotifyMessageReq{
			InvestAmount:       constant.DefaultInvest,
			ExchangeBuy:        Rybit,
			ExchangeSell:       MAX,
			BuyPrice:           qInfo[Rybit].buyPrice,
			SellPrice:          qInfo[MAX].sellPrice,
			Spread:             aInfo.Spread,
			Arbitrage:          aInfo.Arbitrage,
			Profit:             aInfo.Profit,
			IsExcitedArbitrage: isExcitedArbitrage,
			IsExcitedSpread:    isExcitedSpread,
		}); err != nil {
			log.Println(err.Error())
		}
	}
}

type quoteInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Exchanges []struct {
			Name       string          `json:"name"`
			BuyRate    decimal.Decimal `json:"buy_rate"`
			SellRate   decimal.Decimal `json:"sell_rate"`
			UpdateTime int64           `json:"update_time"`
		} `json:"exchanges"`
	} `json:"data"`
}

type quoteResult struct {
	exchange   exchange
	buyPrice   decimal.Decimal
	sellPrice  decimal.Decimal
	updateTime time.Time
}

func fetchExchangeQuotation(exchanges []exchange) (map[exchange]quoteResult, error) {
	httpResp, err := defaultCli.Get("https://www.usdtwhere.com/wallet-api/v1/kgi/exchange-rates/comparison/")
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	qInfo := quoteInfo{}
	err = json.Unmarshal(data, &qInfo)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	result := make(map[exchange]quoteResult)

	fmt.Println(strings.Repeat("=", 20))
	for _, v := range qInfo.Data.Exchanges {
		updateTime := time.UnixMilli(v.UpdateTime)
		fmt.Printf("name: %s, buy: %s, sell: %s, update: %s\n", v.Name, v.BuyRate, v.SellRate, updateTime)
		for _, e := range exchanges {
			if strings.EqualFold(v.Name, string(e)) {
				result[e] = quoteResult{
					exchange:   e,
					buyPrice:   v.BuyRate,
					sellPrice:  v.SellRate,
					updateTime: updateTime,
				}
			}
		}
	}

	return result, nil
}

type arbitrageInfo struct {
	Profit    decimal.Decimal
	Spread    decimal.Decimal
	Arbitrage decimal.Decimal
}

func calArbitrageInfo(invest, buyPrice, sellPrice decimal.Decimal) arbitrageInfo {
	arbitrage := sellPrice.Sub(buyPrice).Div(buyPrice)
	profit := arbitrage.Mul(invest)
	spread := sellPrice.Sub(buyPrice)
	return arbitrageInfo{
		Profit:    profit,
		Arbitrage: arbitrage,
		Spread:    spread,
	}
}

type botSendNotifyMessageReq struct {
	InvestAmount                        decimal.Decimal
	ExchangeBuy                         exchange
	ExchangeSell                        exchange
	BuyPrice                            decimal.Decimal
	SellPrice                           decimal.Decimal
	Spread                              decimal.Decimal
	Arbitrage                           decimal.Decimal
	Profit                              decimal.Decimal
	IsExcitedArbitrage, IsExcitedSpread bool
}

type sendMessageBody struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func botSendNotifyMessage(req botSendNotifyMessageReq) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	arbitrage := req.Arbitrage.Mul(decimal.New(1, 2)).Truncate(2).String() + "%"
	if req.IsExcitedArbitrage {
		arbitrage = fmt.Sprintf("%s%s%s", celebrationEmoji, arbitrage, celebrationEmoji)
	}
	spread := req.Spread.String()
	if req.IsExcitedSpread {
		spread = fmt.Sprintf("%s%s%s", celebrationEmoji, spread, celebrationEmoji)
	}

	reqBody := sendMessageBody{
		ChatID: constant.BotTestGroupChatID,
		Text: fmt.Sprintf(
			msgHtml,
			spread,
			req.InvestAmount,
			req.ExchangeBuy,
			req.BuyPrice,
			req.ExchangeSell,
			req.SellPrice,
			arbitrage,
			req.Profit.Truncate(0),
			constant.BotPersonalChatID,
			constant.Author,
		),
		ParseMode: "HTML",
	}

	data, err := json.Marshal(&reqBody)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = defaultCli.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func errorNotify(errMsg string) {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	reqBody := sendMessageBody{
		ChatID: constant.BotPersonalChatID,
		Text:   "error occur: " + errMsg,
	}

	data, err := json.Marshal(&reqBody)
	if err != nil {
		log.Println(err.Error())
		return
	}

	_, err = defaultCli.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err.Error())
	}
}

var msgHtml = `<strong>&#128060;&#128060;&#128060;  Notify &#128060;&#128060;&#128060;</strong>
<strong>=======================</strong>
<strong>Spread: </strong><u>%s</u>
<strong>Invested Amount: </strong><u>%s</u>
<strong>%s Buy: </strong><u>%s</u>
<strong>%s Sell: </strong><u>%s</u>
<strong>Arbitrage: </strong><u>%s</u>
<strong>Estimated Profit: </strong><u>%s</u>
<strong>Author: </strong><a href="tg://user?id=%s">%s</a>
`

func getUpdateMessage() (*domain.UpdateMessageResponse, error) {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", token, latestUpdateID)

	resp, err := defaultCli.Get(url)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	updateResp := domain.UpdateMessageResponse{}

	err = json.Unmarshal(data, &updateResp)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !updateResp.Ok {
		return nil, err
	}

	return &updateResp, nil
}

func commandFactory(c constant.CommandType) domain.CommandHandler {
	switch c {
	case constant.Alive:
		return command.NewAliveCommand()
	case constant.Help:
		return command.NewHelpCommand()
	case constant.Depth:
		return command.NewDepthCommand()
	default:
		return command.NewUnknownCommand()
	}
}

type botSendMessageReq struct {
	chatID int64
	msg    string
}

func botSendMessage(req botSendMessageReq) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	reqBody := sendMessageBody{
		ChatID: req.chatID,
		Text:   req.msg,
	}

	data, err := json.Marshal(&reqBody)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = defaultCli.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}
