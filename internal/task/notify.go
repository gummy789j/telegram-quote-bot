package task

import (
	"context"
	"log"
	"time"

	"github.com/gummy789j/telegram-quote-bot/internal/config"
	"github.com/gummy789j/telegram-quote-bot/internal/constant"
	domain "github.com/gummy789j/telegram-quote-bot/internal/domain/usecase"
)

type notifyTask struct {
	cfg *config.TelegramCfg
	tb  domain.TelegramUseCase
}

func NewNotifyTask(cfg *config.TelegramCfg, tb domain.TelegramUseCase) Task {
	return &notifyTask{cfg: cfg, tb: tb}
}

func (t *notifyTask) Name() string {
	return "notify"
}

//NOTE: expired at 2025/12/31 23:59:59
func (t *notifyTask) Freq() (runTime time.Duration, tickTime time.Duration) {
	expiredAt := time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC)
	return time.Until(expiredAt), time.Minute
}

func (t *notifyTask) Run(ctx context.Context) error {

	exchangeBuy := constant.Rybit
	exchangeSell := constant.MAX

	toChatID := t.cfg.QuoteComparisonBot.GroupChatID
	if config.IsDevelopment() {
		toChatID = t.cfg.QuoteComparisonBot.TestGroupChatID
	}

	err := t.tb.NotifyArbitrage(ctx, domain.NotifyArbitrageRequest{
		ExchangeBuy:  exchangeBuy,
		ExchangeSell: exchangeSell,
		ToChatID:     toChatID,
	})
	if err != nil {
		log.Println("notify arbitrage job failed: ", err.Error())
	}

	return nil
}
