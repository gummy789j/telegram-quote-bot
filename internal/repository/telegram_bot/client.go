package telegram_bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gummy789j/telegram-quote-bot/internal/config"
	"github.com/gummy789j/telegram-quote-bot/internal/constant"
	domain "github.com/gummy789j/telegram-quote-bot/internal/domain/repo"
	"github.com/gummy789j/telegram-quote-bot/internal/transport"
	"github.com/shopspring/decimal"
)

type telegramBotRepo struct {
	cfg      *config.TelegramCfg
	cli      transport.HttpClient
	endpoint string
}

var _ domain.TelegramBotRepo = (*telegramBotRepo)(nil)

func NewTelegramBotRepo(cli transport.HttpClient, cfg *config.TelegramCfg) domain.TelegramBotRepo {
	return &telegramBotRepo{
		cli:      cli,
		endpoint: fmt.Sprintf("https://api.telegram.org/bot%s", cfg.QuoteComparisonBot.Token),
		cfg:      cfg,
	}
}

var (
	pathSendMessage = "/sendMessage"
	pathGetUpdates  = "/getUpdates"
)

func (t *telegramBotRepo) SendArbitrageNotify(ctx context.Context, req domain.SendArbitrageNotifyRequest) error {

	arbitrage := req.Arbitrage.Mul(decimal.New(1, 2)).Truncate(2).String() + "%"
	if req.IsExcitedArbitrage {
		arbitrage = fmt.Sprintf("%s%s%s", constant.EmojiCelebration, arbitrage, constant.EmojiCelebration)
	}
	spread := req.Spread.String()
	if req.IsExcitedSpread {
		spread = fmt.Sprintf("%s%s%s", constant.EmojiCelebration, spread, constant.EmojiCelebration)
	}

	profit := req.Profit.Truncate(0).String()
	tmpl := tmplArbitrageNotify
	text := tmpl.Format(
		spread,
		req.InvestAmount,
		req.ExchangeBuy,
		req.BuyPrice,
		req.ExchangeSell,
		req.SellPrice,
		arbitrage,
		profit,
		t.cfg.AuthorID,
		t.cfg.Author,
	)

	return t.SendMessage(ctx, domain.SendMessageRequest{
		ChatID:    req.ChatID,
		Text:      text,
		ParseMode: tmpl.Type().String(),
	})
}

func (t *telegramBotRepo) SendErrorNotify(ctx context.Context, req domain.SendErrorNotifyRequest) error {

	tmpl := tmplErrorNotify
	text := tmpl.Format(
		req.Title,
		req.ErrMsg,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return t.SendMessage(ctx, domain.SendMessageRequest{
		ChatID:    req.ChatID,
		Text:      text,
		ParseMode: tmpl.Type().String(),
	})
}

func (t *telegramBotRepo) SendMessage(ctx context.Context, req domain.SendMessageRequest) error {

	url := fmt.Sprintf("%s%s", t.endpoint, pathSendMessage)
	reqBody := map[string]interface{}{
		"chat_id": req.ChatID,
		"text":    req.Text,
	}

	if len(req.ParseMode) > 0 {
		reqBody["parse_mode"] = req.ParseMode
	}

	data, err := json.Marshal(&reqBody)
	if err != nil {
		log.Println("json marshal failed", err.Error())
		return err
	}

	_, err = t.cli.Send(ctx, &transport.HttpRequest{
		Method: http.MethodPost,
		URL:    url,
		Body:   data,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
	if err != nil {
		log.Println("send message failed", err.Error())
		return err
	}

	return nil
}

func (t *telegramBotRepo) GetUpdates(ctx context.Context, req domain.GetUpdatesRequest) (*domain.GetUpdatesResponse, error) {
	url := fmt.Sprintf("%s%s", t.endpoint, pathGetUpdates)

	params := map[string]string{}
	if req.Offset > 0 {
		params["offset"] = fmt.Sprintf("%d", req.Offset)
	}

	httpResp, err := t.cli.Send(ctx, &transport.HttpRequest{
		Method: http.MethodGet,
		URL:    url,
		Params: params,
	})
	if err != nil {
		log.Println("get updates failed", err.Error())
		return nil, err
	}

	resp := &updateMessageResp{}
	if err := json.Unmarshal(httpResp.Body, resp); err != nil {
		log.Println("json unmarshal failed", err.Error())
		return nil, err
	}

	if !resp.Ok {
		log.Println("get updates response nok failed")
		return nil, fmt.Errorf("get updates failed: not ok")
	}

	return &domain.GetUpdatesResponse{
		Result: resp.Result,
	}, nil
}

func (t *telegramBotRepo) GetBotCommandUpdates(ctx context.Context, req domain.GetBotCommandUpdatesRequest) (*domain.GetBotCommandUpdatesResponse, error) {

	getUpdatesResp, err := t.GetUpdates(ctx, domain.GetUpdatesRequest{Offset: req.Offset})
	if err != nil {
		log.Println("get updates failed", err.Error())
		return nil, err
	}

	infos := []*domain.BotCommandInfo{}

	var lastUpdateID *int64 = nil

	for _, v := range getUpdatesResp.Result {

		if lastUpdateID == nil {
			lastUpdateID = &v.UpdateID
		} else {
			if v.UpdateID >= *lastUpdateID {
				*lastUpdateID = v.UpdateID
			}
		}

		if v.Message == nil {
			continue
		}

		if v.Message.Text == nil {
			continue
		}

		cmd := constant.CommandType(strings.TrimSuffix(strings.TrimPrefix(*v.Message.Text, "/"), t.cfg.QuoteComparisonBot.Name))

		if len(v.Message.Entities) == 0 {
			continue
		}

		if v.Message.Entities[0].Type != "bot_command" {
			continue
		}

		if v.Message.From == nil {
			continue
		}

		if v.Message.Chat == nil {
			continue
		}

		infos = append(infos, &domain.BotCommandInfo{
			UpdateID:   v.UpdateID,
			MessageID:  v.Message.MessageID,
			FromChatID: v.Message.Chat.ID,
			FromID:     v.Message.From.ID,
			Command:    cmd,
			Date:       v.Message.Date,
		})
	}

	return &domain.GetBotCommandUpdatesResponse{
		LastUpdateID: lastUpdateID,
		Infos:        infos,
	}, nil
}
