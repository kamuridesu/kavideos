package fetcher

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
)

const (
	USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64; rv:145.0) Gecko/20100101 Firefox/145.0"
)

type CookieFetcher struct {
	Cookies       []*http.Cookie
	RefererHeader string
}

func GenerateRandomFilename(extension string) string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "download" + extension
	}
	return hex.EncodeToString(bytes) + extension
}

func Fetch(ctx context.Context, url string, dst io.Writer, progressCallback func(int, int), cookies ...*CookieFetcher) error {
	slog.Info("Fetching data from url: " + url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", USER_AGENT)

	for _, cookie := range cookies {
		for _, c := range cookie.Cookies {
			req.AddCookie(c)
		}
		if cookie.RefererHeader != "" {
			req.Header.Set("Referer", cookie.RefererHeader)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	for _, cookie := range cookies {
		cookie.Cookies = res.Cookies()
		if ref := res.Header.Get("Referer"); ref != "" {
			cookie.RefererHeader = ref
		}
	}

	if w, ok := dst.(http.ResponseWriter); ok {
		contentType := res.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", contentType)

		if res.ContentLength > 0 {
			w.Header().Set("Content-Length", strconv.Itoa(int(res.ContentLength)))
		}

		ext := filepath.Ext(url)
		if ext == "" {
			exts, _ := mime.ExtensionsByType(contentType)
			if len(exts) > 0 {
				ext = exts[0]
			} else {
				ext = ".mp4"
			}
		}

		filename := GenerateRandomFilename(ext)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		w.WriteHeader(http.StatusOK)
	}

	total := res.ContentLength
	var current int64 = 0

	buf := make([]byte, 32*1024)

	for {
		n, err := res.Body.Read(buf)
		if n > 0 {
			_, wErr := dst.Write(buf[:n])
			if wErr != nil {
				return wErr
			}

			current += int64(n)
			if progressCallback != nil {
				progressCallback(int(total), int(current))
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
