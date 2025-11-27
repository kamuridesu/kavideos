package parser

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseUrlFromHtmlContent(body []byte) (string, error) {
	slog.Info("Parsing HTML body")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	selection := doc.Find("a-video-player").First()

	src, exists := selection.Attr("src")
	if exists {
		return src, nil
	}

	return "", fmt.Errorf("no url found")

}
