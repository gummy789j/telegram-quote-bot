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
	"time"

	"github.com/shopspring/decimal"
)

var defaultCli *http.Client
var token string

var (
	defaultInvest    = decimal.NewFromFloat(500000)
	minSpread        = decimal.NewFromFloat(0.15)
	minArbitrage     = decimal.NewFromFloat(0.005)
	excitedSpread    = decimal.NewFromFloat(0.3)
	excitedArbitrage = decimal.NewFromFloat(0.01)
	notifyFreq       = time.Minute
)

var (
	botGroupChatID    = "-781207517"
	botPersonalChatID = "1881712391"
	author            = "t.me/gummy789j"
)

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
	token = os.Getenv("TELEGRAM_BOT_TOKEN")
	if len(token) == 0 {
		panic("token is empty")
	}
}

func main() {
	var err error

	defer func() {
		errorNotify(err.Error())
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 365*24*time.Hour)
	defer cancel()

	ticker := time.NewTicker(notifyFreq)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println(strings.Repeat("=", 20))
			fmt.Println("work done")
			fmt.Println(strings.Repeat("=", 20))
			return
		case <-ticker.C:
			fmt.Println("time:", time.Now())
			qInfo, err := fetchExchangeQuotation([]exchange{MAX, Rybit})
			if err != nil {
				log.Println(err.Error())
				break
			}

			aInfo := calArbitrageInfo(defaultInvest, qInfo[Rybit].buyPrice, qInfo[MAX].sellPrice)

			if aInfo.Arbitrage.LessThan(minArbitrage) {
				break
			}

			if aInfo.Spread.LessThan(minSpread) {
				break
			}

			var isExcitedArbitrage, isExcitedSpread bool

			if aInfo.Arbitrage.GreaterThanOrEqual(excitedArbitrage) {
				isExcitedArbitrage = true
			}

			if aInfo.Spread.GreaterThanOrEqual(excitedSpread) {
				isExcitedSpread = true
			}

			if err = botSendMessage(botSendMessageReq{
				InvestAmount:       defaultInvest,
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

type botSendMessageReq struct {
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
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func botSendMessage(req botSendMessageReq) error {
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
		ChatID: botGroupChatID,
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
			botPersonalChatID,
			author,
		),
		ParseMode: "HTML",
	}

	data, err := json.Marshal(&reqBody)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	httpResp, err := defaultCli.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func errorNotify(errMsg string) {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	reqBody := sendMessageBody{
		ChatID: botPersonalChatID,
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
