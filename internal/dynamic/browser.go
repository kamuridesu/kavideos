package dynamic

import (
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kamuridesu/go-kwai-dowloader-bot/internal/static"
	"github.com/playwright-community/playwright-go"
)

func FetchAllUrlsInPage(ctx context.Context, rootUrl string, extensions []string) ([]string, error) {
	if extensions == nil {
		extensions = static.DEFAULT_VIDEO_EXT_LIST
	}

	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		return nil, err
	}
	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}
	if _, err := page.Goto(rootUrl); err != nil {
		return nil, err
	}
	content, err := page.Content()
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	var urls []string

	doc.Find("[href], [src]").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			if static.HasExtension(extensions, href) {
				urls = append(urls, href)
			}
		}

		if src, exists := s.Attr("src"); exists {
			if static.HasExtension(extensions, src) {
				urls = append(urls, src)
			}
		}
	})

	return urls, nil

}
