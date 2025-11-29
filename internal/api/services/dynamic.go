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
	cookies := &fetcher.CookieFetcher{RefererHeader: url}
	urls, err := dynamic.FetchAllUrlsInPage(ctx, url, nil, cookies)
	if err != nil {
		return err
	}
	if len(urls) < 1 {
		return fmt.Errorf("no video url found")
	}
	if len(urls) > 1 {
		return downloadAndZipResponse(ctx, urls, w, cookies)
	}

	err = fetcher.Fetch(ctx, urls[0], w, nil, cookies)
	if err != nil {
		return err
	}
	return nil

}
