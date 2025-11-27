package bot

import (
	"bytes"
	"context"
	"log/slog"
	"regexp"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kamuridesu/kavideos/internal/cobalt"
	"github.com/kamuridesu/kavideos/internal/fetcher"
	"github.com/kamuridesu/kavideos/internal/parser"
)

func errorHandler(ctx context.Context, b *bot.Bot, msgId int64, err error) {

	msg := "An error happened while processing request: " + err.Error()
	b.SendMessage(ctx, &bot.SendMessageParams{
		Text:   msg,
		ChatID: msgId,
	})

}

func kwaiHandler(ctx context.Context, kUrl string, b *bot.Bot, update *models.Update) {
	errHandler := func(err error) {
		errorHandler(ctx, b, update.Message.Chat.ID, err)
	}

	html := new(bytes.Buffer)
	err := fetcher.Fetch(ctx, kUrl, html, nil)
	if err != nil {
		errHandler(err)
		return
	}

	url, err := parser.ParseUrlFromHtmlContent(html.Bytes())
	if err != nil {
		errHandler(err)
		return
	}

	prog := progressStart(ctx, b, update.Message.Chat.ID)
	data := new(bytes.Buffer)
	err = fetcher.Fetch(ctx, url, data, prog.update)
	if err != nil {
		errHandler(err)
		return
	}
	prog.finishing()

	params := &bot.SendVideoParams{
		ChatID: update.Message.Chat.ID,
		Video:  &models.InputFileUpload{Filename: "video.mp4", Data: bytes.NewReader(data.Bytes())},
	}
	slog.Info("Sending media")
	b.SendVideo(ctx, params)
	slog.Info("Done")
	prog.end()

}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	txt := update.Message.Text
	slog.Info("New message: " + txt)

	if txt == "" {
		return
	}
	re := regexp.MustCompile(`https://k\.kwai\.com/p/.{8}`)
	kUrl := re.FindString(txt)
	if kUrl != "" {
		kwaiHandler(ctx, kUrl, b, update)
		return
	}

	errHandler := func(err error) {
		errorHandler(ctx, b, update.Message.Chat.ID, err)
	}
	pg := progressStart(ctx, b, update.Message.Chat.ID)
	res := new(bytes.Buffer)
	err := cobalt.DownloadMediaCobalt(ctx, txt, res, pg.update)
	if err != nil {
		errHandler(err)
		return
	}
	pg.finishing()
	params := &bot.SendVideoParams{
		ChatID: update.Message.Chat.ID,
		Video:  &models.InputFileUpload{Filename: "media.mp4", Data: bytes.NewReader(res.Bytes())},
	}
	slog.Info("Sending media")
	b.SendVideo(ctx, params)
	slog.Info("Done")
	pg.end()

}
