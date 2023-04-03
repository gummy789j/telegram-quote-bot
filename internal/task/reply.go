package task

import (
	"context"
	"log"
	"time"

	"github.com/gummy789j/telegram-quote-bot/internal/config"
	domain "github.com/gummy789j/telegram-quote-bot/internal/domain/usecase"
)

type replyTask struct {
	cfg *config.TelegramCfg
	tb  domain.TelegramUseCase
}

func NewReplyTask(cfg *config.TelegramCfg, tb domain.TelegramUseCase) Task {
	return &replyTask{cfg: cfg, tb: tb}
}

func (t *replyTask) Name() string {
	return "notify"
}

//NOTE: expired at 2025/12/31 23:59:59
func (t *replyTask) Freq() (runTime time.Duration, tickTime time.Duration) {
	expiredAt := time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC)
	return time.Until(expiredAt), 2 * time.Second
}

func (t *replyTask) Run(ctx context.Context) error {

	FromChatIDs := []int64{
		t.cfg.QuoteComparisonBot.GroupChatID,
		t.cfg.QuoteComparisonBot.TestGroupChatID,
		t.cfg.AdminChatID,
	}

	if config.IsDevelopment() {
		FromChatIDs = []int64{t.cfg.QuoteComparisonBot.TestGroupChatID, t.cfg.AdminChatID}
	}

	err := t.tb.ReplyCommand(ctx, domain.ReplyCommandRequest{FromChatIDs: FromChatIDs})
	if err != nil {
		log.Println("reply command job failed: ", err.Error())
	}

	return nil
}
