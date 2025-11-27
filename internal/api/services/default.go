package services

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"github.com/kamuridesu/go-kwai-dowloader-bot/internal/cobalt"
)

func DefaultDownloader(ctx context.Context, url string, w io.Writer) error {
	if url == "" {
		return fmt.Errorf("invalid url: empty")
	}
	re := regexp.MustCompile(`https://k\.kwai\.com/p/.{8}`)
	kUrl := re.FindString(url)
	if kUrl != "" {
		return kwaiHandler(ctx, kUrl, w)
	}

	err := cobalt.DownloadMediaCobalt(ctx, url, w, nil)
	if err != nil {
		return err
	}
	return nil
}
