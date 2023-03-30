package command

import (
	"fmt"

	"github.com/gummy789j/telegram-quote-bot/internal/constant"
	"github.com/gummy789j/telegram-quote-bot/internal/domain"
	"github.com/shopspring/decimal"
)

type aliveCommand struct {
}

func NewAliveCommand() domain.CommandHandler {
	return &aliveCommand{}
}

func (c *aliveCommand) Reply(id int64) string {
	if id == constant.AuthorID {
		return "I'm alive"
	} else {
		return "我是懶惰老鼠，但不是死老鼠"
	}
}

type helpCommand struct {
}

func NewHelpCommand() domain.CommandHandler {
	return &helpCommand{}
}

func (c *helpCommand) Reply(id int64) string {
	return fmt.Sprintf("我是懶惰老鼠，只喜歡搬 spread 大於 %s 而且 arbitrage 大於 %s%%的單", constant.MinSpread, constant.MinArbitrage.Mul(decimal.New(1, 2)))
}

type depthCommand struct {
}

func NewDepthCommand() domain.CommandHandler {
	return &depthCommand{}
}

func (c *depthCommand) Reply(id int64) string {
	return "我是懶惰老鼠，還沒串這個功能"
}

type unknownCommand struct {
}

func NewUnknownCommand() domain.CommandHandler {
	return &unknownCommand{}
}

func (c *unknownCommand) Reply(id int64) string {
	return "我是懶惰老鼠，不知道你在說什麼"
}
