package services

import (
	"context"
	"fmt"
	"io"

	"github.com/kamuridesu/kavideos/internal/dynamic"
	"github.com/kamuridesu/kavideos/internal/fetcher"
)

func DownloadFromDynamic(ctx context.Context, url string, w io.Writer) error {
	if url == "" {
		return fmt.Errorf("invalid url: empty")
	}
	urls, err := dynamic.FetchAllUrlsInPage(ctx, url, nil)
	if err != nil {
		return err
	}
	if len(urls) < 1 {
		return fmt.Errorf("no video url found")
	}
	if len(urls) > 1 {
		return downloadAndZipResponse(ctx, urls, w)
	}

	err = fetcher.Fetch(ctx, urls[0], w, nil)
	if err != nil {
		return err
	}
	return nil

}
