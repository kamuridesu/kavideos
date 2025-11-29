package static

import (
	"bytes"
	"context"
	"log/slog"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kamuridesu/kavideos/internal/fetcher"
)

func HasExtension(extensions []string, url string) bool {
	url = strings.Split(url, "?")[0]
	url = strings.Split(url, "#")[0]
	for _, ext := range extensions {
		if strings.HasSuffix(url, "."+ext) {
			return true
		}
	}
	return false
}

// Passing nil to the extensions list will cause it to load static.DEFAULT_VIDEO_EXT_LIST
func FetchAllUrlsInPage(ctx context.Context, rootUrl string, extensions []string, cookies ...*fetcher.CookieFetcher) ([]string, error) {
	slog.Info("Fetching page urls")

	if extensions == nil {
		extensions = DEFAULT_VIDEO_EXT_LIST
	}

	buffer := new(bytes.Buffer)
	err := fetcher.Fetch(ctx, rootUrl, buffer, nil, cookies...)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(buffer.String()))
	if err != nil {
		return nil, err
	}

	var urls []string

	doc.Find("[href], [src]").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			if HasExtension(extensions, href) && !slices.Contains(urls, href) {
				urls = append(urls, href)
			}
		}

		if src, exists := s.Attr("src"); exists {
			if HasExtension(extensions, src) && !slices.Contains(urls, src) {
				urls = append(urls, src)
			}
		}
	})

	return urls, nil

}
