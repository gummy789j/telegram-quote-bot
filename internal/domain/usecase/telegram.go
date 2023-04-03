package domain

import (
	"context"

	"github.com/gummy789j/telegram-quote-bot/internal/constant"
)

type TelegramUseCase interface {
	ReplyCommand(ctx context.Context, req ReplyCommandRequest) error
	NotifyArbitrage(ctx context.Context, req NotifyArbitrageRequest) error
}

type ReplyCommandRequest struct {
	FromChatIDs []int64
}

func (r *ReplyCommandRequest) IsEmpty() bool {
	return len(r.FromChatIDs) == 0
}

func (r *ReplyCommandRequest) InFromChatIDs(chatID int64) bool {
	for _, v := range r.FromChatIDs {
		if v == chatID {
			return true
		}
	}
	return false
}

type NotifyArbitrageRequest struct {
	ExchangeBuy  constant.Exchange
	ExchangeSell constant.Exchange
	ToChatID     int64
}
