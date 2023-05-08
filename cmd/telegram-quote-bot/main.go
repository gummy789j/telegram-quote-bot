package main

import (
	"context"
	"log"
	"time"

	tqb "github.com/gummy789j/telegram-quote-bot/internal/cmd/telegram-quote-bot"
	"github.com/gummy789j/telegram-quote-bot/internal/config"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()

	ge := tqb.RunServer(ctx, cfg)

	if err := ge.Run(":" + cfg.APIServer.Port); err != nil {
		log.Println("run server failed: ", err.Error())
		return
	}

	log.Println("run telegram server done")

	<-ctx.Done()

	// graceful shutdown
	time.Sleep(10 * time.Second)
}
