package bot

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start(ctx context.Context) error {

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		return err
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startAndHelpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, startAndHelpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "browser", bot.MatchTypeCommand, browserHandler)

	slog.Info("Bot started")
	b.Start(ctx)
	return nil
}

func startAndHelpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "You can send me a social media video link and I will download it for you. \nYou can also use /download to download all videos from any webpage",
	})
}
