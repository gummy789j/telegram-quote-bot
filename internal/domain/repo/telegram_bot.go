package domain

import (
	"context"

	"github.com/gummy789j/telegram-quote-bot/internal/constant"
	"github.com/shopspring/decimal"
)

type TelegramBotRepo interface {
	SendMessage(ctx context.Context, req SendMessageRequest) error
	SendArbitrageNotify(ctx context.Context, req SendArbitrageNotifyRequest) error
	SendErrorNotify(ctx context.Context, req SendErrorNotifyRequest) error
	GetUpdates(ctx context.Context, req GetUpdatesRequest) (*GetUpdatesResponse, error)
	GetBotCommandUpdates(ctx context.Context, req GetBotCommandUpdatesRequest) (*GetBotCommandUpdatesResponse, error)
}

type SendMessageRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type SendArbitrageNotifyRequest struct {
	ChatID                              int64
	InvestAmount                        decimal.Decimal
	ExchangeBuy                         constant.Exchange
	ExchangeSell                        constant.Exchange
	BuyPrice                            decimal.Decimal
	SellPrice                           decimal.Decimal
	Spread                              decimal.Decimal
	Arbitrage                           decimal.Decimal
	Profit                              decimal.Decimal
	IsExcitedArbitrage, IsExcitedSpread bool
}

type SendErrorNotifyRequest struct {
	ChatID int64
	Title  string
	ErrMsg string
}

type GetUpdatesRequest struct {
	Offset int64 `json:"offset"`
}

type GetUpdatesResponse struct {
	Result []struct {
		UpdateID int64 `json:"update_id"`
		Message  *struct {
			MessageID int64 `json:"message_id"`
			From      *struct {
				ID        int64  `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"from"`
			Chat *struct {
				ID                          int64  `json:"id"`
				Title                       string `json:"title"`
				Type                        string `json:"type"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			} `json:"chat"`
			Date               int64 `json:"date"`
			NewChatParticipant *struct {
				ID        int64  `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"new_chat_participant"`
			NewChatMember *struct {
				ID        int64  `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"new_chat_member"`
			NewChatMembers []struct {
				ID        int64  `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"new_chat_members"`
			Text     *string `json:"text"`
			Entities []struct {
				Offset int64  `json:"offset"`
				Length int64  `json:"length"`
				Type   string `json:"type"`
			} `json:"entities"`
		} `json:"message"`
	} `json:"result"`
}

type GetBotCommandUpdatesRequest struct {
	Offset int64 `json:"offset"`
}

type GetBotCommandUpdatesResponse struct {
	LastUpdateID *int64
	Infos        []*BotCommandInfo
}

type BotCommandInfo struct {
	UpdateID   int64
	MessageID  int64
	FromChatID int64
	FromID     int64
	Command    constant.CommandType
	Date       int64
}
