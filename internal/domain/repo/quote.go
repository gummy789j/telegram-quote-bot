package domain

import (
	"context"
	"time"

	"github.com/gummy789j/telegram-quote-bot/internal/constant"
	"github.com/shopspring/decimal"
)

type QuoteRepo interface {
	GetQuotations(ctx context.Context, req GetQuotationsRequest) (*GetQuotationsResponse, error)
}

type GetQuotationsRequest struct {
}

type GetQuotationsResponse struct {
	Infos map[constant.Exchange]QuotationInfo
}

type QuotationInfo struct {
	BuyPrice   decimal.Decimal
	SellPrice  decimal.Decimal
	UpdateTime time.Time
}
