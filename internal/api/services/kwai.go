package services

import (
	"bytes"
	"context"
	"io"

	"github.com/kamuridesu/kavideos/internal/fetcher"
	"github.com/kamuridesu/kavideos/internal/parser"
)

func kwaiHandler(ctx context.Context, url string, w io.Writer) error {
	html := new(bytes.Buffer)
	cookie := &fetcher.CookieFetcher{RefererHeader: url}
	err := fetcher.Fetch(ctx, url, html, nil, cookie)
	if err != nil {
		return err
	}
	vUrl, err := parser.ParseUrlFromHtmlContent(html.Bytes())
	if err != nil {
		return err
	}

	err = fetcher.Fetch(ctx, vUrl, w, nil, cookie)
	if err != nil {
		return err
	}
	return nil
}
