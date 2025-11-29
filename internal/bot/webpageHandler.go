package bot

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	nurl "net/url"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kamuridesu/kavideos/internal/dynamic"
	"github.com/kamuridesu/kavideos/internal/fetcher"
)

func downloadAndSend(ctx context.Context, b *bot.Bot, update *models.Update, url string, cookie ...*fetcher.CookieFetcher) {
	errHandler := func(err error) {
		errorHandler(ctx, b, update.Message.Chat.ID, err)
	}

	pg := progressStart(ctx, b, update.Message.Chat.ID)

	data := new(bytes.Buffer)
	err := fetcher.Fetch(ctx, url, data, pg.update, cookie...)
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

func browserHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	errHandler := func(err error) {
		errorHandler(ctx, b, update.Message.Chat.ID, err)
	}

	url := strings.TrimSpace(strings.TrimLeft(update.Message.Text, "/browser"))
	_, err := nurl.ParseRequestURI(url)
	if url == "" || !strings.HasPrefix(url, "http") || err != nil {
		errHandler(fmt.Errorf("invalid url"))
		return
	}

	cookie := &fetcher.CookieFetcher{RefererHeader: url}
	urls, err := dynamic.FetchAllUrlsInPage(ctx, url, nil, cookie)
	if err != nil {
		errHandler(err)
		return
	}

	for _, url := range urls {
		go downloadAndSend(ctx, b, update, url, cookie)
	}

}
