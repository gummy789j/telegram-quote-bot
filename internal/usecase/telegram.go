package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gummy789j/telegram-quote-bot/internal/config"
	"github.com/gummy789j/telegram-quote-bot/internal/constant"
	dRepo "github.com/gummy789j/telegram-quote-bot/internal/domain/repo"
	dUc "github.com/gummy789j/telegram-quote-bot/internal/domain/usecase"
	"github.com/shopspring/decimal"
)

type telegramUseCase struct {
	cfg   *config.TelegramCfg
	tb    dRepo.TelegramBotRepo
	quote dRepo.QuoteRepo

	// mutex
	lock *sync.Mutex
}

var _ dUc.TelegramUseCase = (*telegramUseCase)(nil)

var latestUpdateID int64

func NewTelegramUseCase(cfg *config.TelegramCfg, tb dRepo.TelegramBotRepo, quote dRepo.QuoteRepo) dUc.TelegramUseCase {
	uc := &telegramUseCase{cfg: cfg, tb: tb, quote: quote, lock: &sync.Mutex{}}

	// get the latest update id and store it
	umResp, err := uc.tb.GetUpdates(context.Background(), dRepo.GetUpdatesRequest{})
	if err != nil {
		panic("get latest update id failed: " + err.Error())
	}

	if len(umResp.Result) != 0 {
		for _, v := range umResp.Result {
			if v.UpdateID > latestUpdateID {
				latestUpdateID = v.UpdateID
			}
		}
	}
	return uc
}

func (u *telegramUseCase) ReplyCommand(ctx context.Context, req dUc.ReplyCommandRequest) error {

	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
		if err != nil {
			u.notifyError(ctx, "ReplyCommand", err.Error())
		}
	}()

	u.lock.Lock()
	defer u.lock.Unlock()

	// get updates
	cuResp, err := u.tb.GetBotCommandUpdates(ctx, dRepo.GetBotCommandUpdatesRequest{
		Offset: latestUpdateID,
	})
	if err != nil {
		log.Println("get bot commands updates failed: ", err.Error())
		return err
	}

	// reply to the command
	for _, v := range cuResp.Infos {
		if v.UpdateID <= latestUpdateID {
			continue
		}

		if !req.InFromChatIDs(v.FromChatID) {
			continue
		}

		if err := newCommandFactory(commandFactoryReq{
			cfg:         u.cfg,
			commandType: v.Command,
			tb:          u.tb,
			quote:       u.quote,
		}).Reply(v.FromID, v.FromChatID); err != nil {
			log.Println("reply command failed: ", err.Error())
			return err
		}
	}

	if cuResp.LastUpdateID != nil {
		latestUpdateID = *cuResp.LastUpdateID
	}

	return nil
}

func (u *telegramUseCase) NotifyArbitrage(ctx context.Context, req dUc.NotifyArbitrageRequest) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
		if err != nil {
			u.notifyError(ctx, "NotifyArbitrage", err.Error())
		}
	}()

	// get comparison quote
	qInfo, err := u.quote.GetQuotations(ctx, dRepo.GetQuotationsRequest{})
	if err != nil {
		log.Println("get quotations failed: ", err.Error())
		return err
	}

	// calculate arbitrage info
	aInfo := calArbitrageInfo(u.cfg.QuoteComparisonBot.DefaultInvest, qInfo.Infos[req.ExchangeBuy].BuyPrice, qInfo.Infos[req.ExchangeSell].SellPrice)

	if aInfo.Arbitrage.LessThan(u.cfg.QuoteComparisonBot.MinArbitrage) ||
		aInfo.Spread.LessThan(u.cfg.QuoteComparisonBot.MinSpread) {
		return nil
	}

	var isExcitedArbitrage, isExcitedSpread bool

	if aInfo.Arbitrage.GreaterThanOrEqual(u.cfg.QuoteComparisonBot.ExcitedArbitrage) {
		isExcitedArbitrage = true
	}

	if aInfo.Spread.GreaterThanOrEqual(u.cfg.QuoteComparisonBot.ExcitedSpread) {
		isExcitedSpread = true
	}

	// send arbitrage notify
	err = u.tb.SendArbitrageNotify(ctx, dRepo.SendArbitrageNotifyRequest{
		ChatID:             req.ToChatID,
		InvestAmount:       u.cfg.QuoteComparisonBot.DefaultInvest,
		ExchangeBuy:        req.ExchangeBuy,
		ExchangeSell:       req.ExchangeSell,
		BuyPrice:           qInfo.Infos[req.ExchangeBuy].BuyPrice,
		SellPrice:          qInfo.Infos[req.ExchangeSell].SellPrice,
		Spread:             aInfo.Spread,
		Arbitrage:          aInfo.Arbitrage,
		Profit:             aInfo.Profit,
		IsExcitedArbitrage: isExcitedArbitrage,
		IsExcitedSpread:    isExcitedSpread,
	})
	if err != nil {
		log.Println("send arbitrage notify failed: ", err.Error())
		return err
	}
	return nil
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

func (u *telegramUseCase) notifyError(ctx context.Context, title string, errMsg string) {
	// send error notify
	err := u.tb.SendErrorNotify(ctx, dRepo.SendErrorNotifyRequest{
		ChatID: u.cfg.AdminChatID,
		Title:  title,
		ErrMsg: errMsg,
	})
	if err != nil {
		log.Println("send error notify failed: ", err.Error())
	}
}

// command handler
type commandFactoryReq struct {
	cfg         *config.TelegramCfg
	commandType constant.CommandType
	tb          dRepo.TelegramBotRepo
	quote       dRepo.QuoteRepo
}

func newCommandFactory(req commandFactoryReq) commandHandler {
	switch req.commandType {
	case constant.Alive:
		return newAliveCommand(req.cfg, req.tb)
	case constant.Help:
		return newHelpCommand(req.cfg, req.tb)
	case constant.Depth:
		return newDepthCommand(req.tb)
	case constant.Arbitrage:
		return newArbitrageCommand(req.cfg, req.tb, req.quote)
	default:
		return newUnknownCommand(req.tb)
	}
}

type commandHandler interface {
	Reply(toID int64, chatID int64) error
}

type aliveCommand struct {
	cfg *config.TelegramCfg
	tb  dRepo.TelegramBotRepo
}

func newAliveCommand(cfg *config.TelegramCfg, tb dRepo.TelegramBotRepo) commandHandler {
	return &aliveCommand{cfg: cfg, tb: tb}
}

func (c *aliveCommand) Reply(toID int64, chatID int64) error {
	var msg string
	if toID == c.cfg.AuthorID {
		msg = "I'm alive"
	} else {
		msg = "我是懶惰老鼠，但不是死老鼠"
	}

	return c.tb.SendMessage(context.Background(), dRepo.SendMessageRequest{
		ChatID: chatID,
		Text:   msg,
	})
}

type helpCommand struct {
	cfg *config.TelegramCfg
	tb  dRepo.TelegramBotRepo
}

func newHelpCommand(cfg *config.TelegramCfg, tb dRepo.TelegramBotRepo) commandHandler {
	return &helpCommand{cfg: cfg, tb: tb}
}

func (c *helpCommand) Reply(toID int64, chatID int64) error {
	msg := fmt.Sprintf("我是懶惰老鼠，只喜歡搬 spread 大於 %s 而且 arbitrage 大於 %s%%的單", c.cfg.QuoteComparisonBot.MinSpread, c.cfg.QuoteComparisonBot.MinArbitrage.Mul(decimal.New(1, 2)))

	return c.tb.SendMessage(context.Background(), dRepo.SendMessageRequest{
		ChatID: chatID,
		Text:   msg,
	})
}

type depthCommand struct {
	tb dRepo.TelegramBotRepo
}

func newDepthCommand(tb dRepo.TelegramBotRepo) commandHandler {
	return &depthCommand{tb: tb}
}

func (c *depthCommand) Reply(toID int64, chatID int64) error {
	msg := "我是懶惰老鼠，還沒串這個功能"

	return c.tb.SendMessage(context.Background(), dRepo.SendMessageRequest{
		ChatID: chatID,
		Text:   msg,
	})
}

type arbitrageCommand struct {
	cfg   *config.TelegramCfg
	quote dRepo.QuoteRepo
	tb    dRepo.TelegramBotRepo
}

func newArbitrageCommand(cfg *config.TelegramCfg, tb dRepo.TelegramBotRepo, quote dRepo.QuoteRepo) commandHandler {
	return &arbitrageCommand{cfg: cfg, tb: tb, quote: quote}
}

func (c *arbitrageCommand) Reply(toID int64, chatID int64) error {

	// send message
	ctx := context.Background()
	// get comparison quote
	qInfo, err := c.quote.GetQuotations(ctx, dRepo.GetQuotationsRequest{})
	if err != nil {
		log.Println("get quotations failed: ", err.Error())
		return err
	}

	// calculate arbitrage info
	aInfo := calArbitrageInfo(c.cfg.QuoteComparisonBot.DefaultInvest, qInfo.Infos[constant.Rybit].BuyPrice, qInfo.Infos[constant.MAX].SellPrice)

	var isExcitedArbitrage, isExcitedSpread bool

	if aInfo.Arbitrage.GreaterThanOrEqual(c.cfg.QuoteComparisonBot.ExcitedArbitrage) {
		isExcitedArbitrage = true
	}

	if aInfo.Spread.GreaterThanOrEqual(c.cfg.QuoteComparisonBot.ExcitedSpread) {
		isExcitedSpread = true
	}

	// send arbitrage notify
	return c.tb.SendArbitrageNotify(ctx, dRepo.SendArbitrageNotifyRequest{
		ChatID:             chatID,
		InvestAmount:       c.cfg.QuoteComparisonBot.DefaultInvest,
		ExchangeBuy:        constant.Rybit,
		ExchangeSell:       constant.MAX,
		BuyPrice:           qInfo.Infos[constant.Rybit].BuyPrice,
		SellPrice:          qInfo.Infos[constant.MAX].SellPrice,
		Spread:             aInfo.Spread,
		Arbitrage:          aInfo.Arbitrage,
		Profit:             aInfo.Profit,
		IsExcitedArbitrage: isExcitedArbitrage,
		IsExcitedSpread:    isExcitedSpread,
	})
}

type unknownCommand struct {
	tb dRepo.TelegramBotRepo
}

func newUnknownCommand(tb dRepo.TelegramBotRepo) commandHandler {
	return &unknownCommand{tb: tb}
}

func (c *unknownCommand) Reply(toID int64, chatID int64) error {
	msg := "我是懶惰老鼠，不知道你在說什麼"
	return c.tb.SendMessage(context.Background(), dRepo.SendMessageRequest{
		ChatID: chatID,
		Text:   msg,
	})
}
