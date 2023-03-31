package constant

import (
	"time"

	"github.com/shopspring/decimal"
)

var (
	BotGroupChatID     int64  = -781207517
	BotTestGroupChatID int64  = -905284654
	BotPersonalChatID  int64  = 1881712391
	AuthorID           int64  = 1881712391
	Author                    = "t.me/gummy789j"
	BotName            string = "@gummy_s_bot"
)

var (
	DefaultInvest    = decimal.NewFromFloat(500000)
	MinSpread        = decimal.NewFromFloat(0.15)
	MinArbitrage     = decimal.NewFromFloat(0.005)
	ExcitedSpread    = decimal.NewFromFloat(0.3)
	ExcitedArbitrage = decimal.NewFromFloat(0.01)
	NotifyFreq       = time.Minute
	ReplyFreq        = 2 * time.Second
)

type CommandType string

var (
	Alive     CommandType = "alive"
	Help      CommandType = "help"
	Depth     CommandType = "depth"
	Arbitrage CommandType = "arbitrage"
)
