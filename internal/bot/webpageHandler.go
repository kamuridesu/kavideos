package bot

import (
	"bytes"
	"context"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kamuridesu/go-kwai-dowloader-bot/internal/dynamic"
	"github.com/kamuridesu/go-kwai-dowloader-bot/internal/fetcher"
	"github.com/kamuridesu/go-kwai-dowloader-bot/internal/static"
)

func downloadAndSend(ctx context.Context, b *bot.Bot, update *models.Update, url string) {
	errHandler := func(err error) {
		errorHandler(ctx, b, update.Message.Chat.ID, err)
	}

	pg := progressStart(ctx, b, update.Message.Chat.ID)

	data := new(bytes.Buffer)
	err := fetcher.Fetch(ctx, url, data, pg.update)
	if err != nil {
		errHandler(err)
		return
	}
	pg.finishing()

	params := &bot.SendVideoParams{
		ChatID: update.Message.Chat.ID,
		Video:  &models.InputFileUpload{Filename: "video.mp4", Data: bytes.NewReader(data.Bytes())},
	}

	b.SendVideo(ctx, params)

	pg.end()
}

func webpageHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info(update.Message.Text)
	errHandler := func(err error) {
		errorHandler(ctx, b, update.Message.Chat.ID, err)
	}

	url := strings.TrimSpace(strings.TrimLeft(update.Message.Text, "/download"))

	urls, err := static.FetchAllUrlsInPage(ctx, url, nil)
	if err != nil {
		errHandler(err)
		return
	}

	for _, url := range urls {
		go downloadAndSend(ctx, b, update, url)
	}
}

func browserHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info(update.Message.Text)
	errHandler := func(err error) {
		errorHandler(ctx, b, update.Message.Chat.ID, err)
	}

	url := strings.TrimSpace(strings.TrimLeft(update.Message.Text, "/browser"))

	urls, err := dynamic.FetchAllUrlsInPage(ctx, url, nil)
	if err != nil {
		errHandler(err)
		return
	}

	for _, url := range urls {
		go downloadAndSend(ctx, b, update, url)
	}

}
