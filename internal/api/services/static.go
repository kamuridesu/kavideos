package services

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/kamuridesu/go-kwai-dowloader-bot/internal/fetcher"
	"github.com/kamuridesu/go-kwai-dowloader-bot/internal/static"
)

type StructuredResponse struct {
	Blob  []byte
	Error error
}

func downloadMedia(ctx context.Context, wg *sync.WaitGroup, ch chan StructuredResponse, url string) {
	defer wg.Done()
	var response StructuredResponse
	data := new(bytes.Buffer)
	err := fetcher.Fetch(ctx, url, data, nil)
	if err != nil {
		response.Error = err
	} else {
		response.Blob = data.Bytes()
	}
	select {
	case ch <- response:
	case <-ctx.Done():
		return
	}
}

func downloadAndZipResponse(ctx context.Context, urls []string, w io.Writer) error {
	results := make(chan StructuredResponse, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go downloadMedia(ctx, &wg, results, url)
	}

	wg.Wait()
	close(results)

	buffer := new(bytes.Buffer)
	zipW := zip.NewWriter(buffer)
	c := 1
	for res := range results {
		if res.Error != nil {
			slog.Error("failed to download media", "err", res.Error)
		}

		if len(res.Blob) == 0 {
			continue
		}

		file, err := zipW.Create(fmt.Sprintf("media-{%d}.mp4", c))
		if err != nil {
			return err
		}
		_, err = file.Write(res.Blob)
		if err != nil {
			continue
		}
		c++
	}
	if c == 1 {
		return fmt.Errorf("all downloads failed")
	}
	err := zipW.Close()

	if err != nil {
		return err
	}

	if w, ok := w.(http.ResponseWriter); ok {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))
		w.Header().Set("Content-Disposition", "attachment; filename="+fetcher.GenerateRandomFilename("zip"))
		w.WriteHeader(http.StatusOK)
	}

	w.Write(buffer.Bytes())
	return nil
}

func DownloadFromStatic(ctx context.Context, url string, w io.Writer) error {
	if url == "" {
		return fmt.Errorf("invalid url: empty")
	}
	urls, err := static.FetchAllUrlsInPage(ctx, url, nil)
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
