package config

import (
	"os"

	"github.com/shopspring/decimal"
)

type Config struct {
	APIServer *APIServerCfg
	Telegram  *TelegramCfg
}

func NewConfig(isDev ...bool) *Config {
	if len(isDev) > 0 {
		isDevelopment = isDev[0]
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if len(telegramToken) == 0 {
		panic("TELEGRAM_BOT_TOKEN is not set")
	}
	return &Config{
		APIServer: &APIServerCfg{
			Port: port,
		},
		Telegram: &TelegramCfg{
			AdminChatID: 1881712391,
			AuthorID:    1881712391,
			Author:      "t.me/gummy789j",
			QuoteComparisonBot: &quoteComparisonBot{
				Token:           telegramToken,
				Name:            "@gummy_s_bot",
				GroupChatID:     -781207517,
				TestGroupChatID: -905284654,

				// info
				DefaultInvest:    decimal.NewFromFloat(500000),
				MinSpread:        decimal.NewFromFloat(0.1),
				MinArbitrage:     decimal.NewFromFloat(0.005),
				ExcitedSpread:    decimal.NewFromFloat(0.3),
				ExcitedArbitrage: decimal.NewFromFloat(0.01),
			},
		},
	}
}

type APIServerCfg struct {
	Port string
}

type TelegramCfg struct {
	AdminChatID        int64
	AuthorID           int64
	Author             string
	QuoteComparisonBot *quoteComparisonBot
}

type quoteComparisonBot struct {
	Token            string
	Name             string
	GroupChatID      int64
	TestGroupChatID  int64
	DefaultInvest    decimal.Decimal
	MinSpread        decimal.Decimal
	MinArbitrage     decimal.Decimal
	ExcitedSpread    decimal.Decimal
	ExcitedArbitrage decimal.Decimal
}

var isDevelopment bool = false

func IsDevelopment() bool {
	return isDevelopment
}
