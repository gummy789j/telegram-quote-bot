package main

import (
	"context"
	"log"
	"sync"

	tqb "github.com/gummy789j/telegram-quote-bot/internal/cmd/telegram-quote-bot"
	"github.com/gummy789j/telegram-quote-bot/internal/config"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()
	wg := &sync.WaitGroup{}

	tqb.RunServer(ctx, cfg, wg)

	log.Println("run telegram server done")

	wg.Wait()
}
