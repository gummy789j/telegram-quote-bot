package server

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gummy789j/telegram-quote-bot/internal/config"
	comp "github.com/gummy789j/telegram-quote-bot/internal/repository/comparison"
	tb "github.com/gummy789j/telegram-quote-bot/internal/repository/telegram_bot"
	"github.com/gummy789j/telegram-quote-bot/internal/task"
	"github.com/gummy789j/telegram-quote-bot/internal/transport"
	"github.com/gummy789j/telegram-quote-bot/internal/usecase"
)

func RunServer(ctx context.Context, cfg *config.Config) *gin.Engine {

	telegramBotRepo := tb.NewTelegramBotRepo(transport.NewHttpClient(), cfg.Telegram)
	comparisonRepo := comp.NewComparisonClient(transport.NewHttpClient())
	telegramUseCase := usecase.NewTelegramUseCase(cfg.Telegram, telegramBotRepo, comparisonRepo)

	tasks := []task.Task{
		task.NewNotifyTask(cfg.Telegram, telegramUseCase),
		task.NewReplyTask(cfg.Telegram, telegramUseCase),
	}

	jobProcessor(ctx, tasks)

	g := gin.New()

	g.GET("/", func(c *gin.Context) {
		c.String(200, "alive")
	})

	g.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return g
}

func jobProcessor(pctx context.Context, tasks []task.Task) {
	for _, t := range tasks {

		go func(pctx context.Context, t task.Task) {

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

		}(pctx, t)
	}
}
