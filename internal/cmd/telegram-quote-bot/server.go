package server

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gummy789j/telegram-quote-bot/internal/config"
	comp "github.com/gummy789j/telegram-quote-bot/internal/repository/comparison"
	tb "github.com/gummy789j/telegram-quote-bot/internal/repository/telegram_bot"
	"github.com/gummy789j/telegram-quote-bot/internal/task"
	"github.com/gummy789j/telegram-quote-bot/internal/transport"
	"github.com/gummy789j/telegram-quote-bot/internal/usecase"
)

func RunServer(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup) {

	telegramBotRepo := tb.NewTelegramBotRepo(transport.NewHttpClient(), cfg.Telegram)
	comparisonRepo := comp.NewComparisonClient(transport.NewHttpClient())
	telegramUseCase := usecase.NewTelegramUseCase(cfg.Telegram, telegramBotRepo, comparisonRepo)

	tasks := []task.Task{
		task.NewNotifyTask(cfg.Telegram, telegramUseCase),
		task.NewReplyTask(cfg.Telegram, telegramUseCase),
	}

	wg.Add(len(tasks))

	jobProcessor(ctx, wg, tasks)
}

func jobProcessor(pctx context.Context, wg *sync.WaitGroup, tasks []task.Task) {
	for _, t := range tasks {

		go func(pctx context.Context, wg *sync.WaitGroup, t task.Task) {
			defer wg.Done()

			runTime, tickTime := t.Freq()

			ctx, cancel := context.WithTimeout(pctx, runTime)
			defer cancel()

			ticker := time.NewTicker(tickTime)
			defer ticker.Stop()

			for {
				t.Run(ctx)
				select {
				case <-ctx.Done():
					log.Println("job done: ", t.Name())
					return

				case <-ticker.C:
				}
			}

		}(pctx, wg, t)
	}
}
