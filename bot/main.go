package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/kamuridesu/go-kwai-dowloader-bot/internal/bot"
)

func main() {

	ctx := context.TODO()
	err := bot.Start(ctx)
	if err != nil {
		slog.Error("An error happened while starting the bot: ", "error", err)
		os.Exit(1)
	}

}
