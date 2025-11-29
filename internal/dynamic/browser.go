package dynamic

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/kamuridesu/kavideos/internal/fetcher"
	"github.com/kamuridesu/kavideos/internal/static"
	"github.com/playwright-community/playwright-go"
)

func ConvertPlaywrightCookies(pwCookies []playwright.Cookie) []*http.Cookie {
	var httpCookies []*http.Cookie

	for _, pwCookie := range pwCookies {
		c := &http.Cookie{
			Name:     pwCookie.Name,
			Value:    pwCookie.Value,
			Domain:   pwCookie.Domain,
			Path:     pwCookie.Path,
			Secure:   pwCookie.Secure,
			HttpOnly: pwCookie.HttpOnly,
		}

		if pwCookie.Expires > 0 {
			c.Expires = time.Unix(int64(pwCookie.Expires), 0)
		}

		switch strings.ToLower(string(*pwCookie.SameSite)) {
		case "lax":
			c.SameSite = http.SameSiteLaxMode
		case "strict":
			c.SameSite = http.SameSiteStrictMode
		case "none":
			c.SameSite = http.SameSiteNoneMode
		default:
			c.SameSite = http.SameSiteDefaultMode
		}

		httpCookies = append(httpCookies, c)
	}

	return httpCookies
}

func FetchAllUrlsInPage(ctx context.Context, rootUrl string, extensions []string, cookieStore ...*fetcher.CookieFetcher) ([]string, error) {
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

	cookies, err := page.Context().Cookies()
	if err != nil {
		return nil, err
	}

	httpCookies := ConvertPlaywrightCookies(cookies)
	for _, c := range cookieStore {
		c.Cookies = httpCookies
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	var urls []string

	doc.Find("[href], [src]").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			if static.HasExtension(extensions, href) && !slices.Contains(urls, href) {
				urls = append(urls, href)
			}
		}

		if src, exists := s.Attr("src"); exists {
			if static.HasExtension(extensions, src) && !slices.Contains(urls, src) {
				urls = append(urls, src)
			}
		}
	})

	return urls, nil

}
