package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/kamuridesu/go-kwai-dowloader-bot/internal/api/services"
	"github.com/kamuridesu/gomechan/core/response"
)

type ExpectedRequest struct {
	Url        string `json:"url"`
	IsWebpage  bool   `json:"isWebpage"`
	UseBrowser bool   `json:"useBrowser"`
}

func errorHandler(res *response.ResponseWriter, err error, statusCode ...int) {
	if errors.Is(err, context.Canceled) {
		return
	}
	slog.Error("an error happened while processing request", "error", err)
	status := 500
	if len(statusCode) > 0 {
		status = statusCode[0]
	}
	res.AsJson(status, map[string]any{"error": err.Error()})
}

func defaultDownload(w http.ResponseWriter, r *http.Request) {
	res := response.New(&w, r)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorHandler(&res, err)
		return
	}
	defer r.Body.Close()
	var expectedRequest ExpectedRequest
	err = json.Unmarshal(body, &expectedRequest)
	if err != nil {
		errorHandler(&res, err)
		return
	}
	if expectedRequest.Url == "" {
		errorHandler(&res, fmt.Errorf("no url found in request"), 400)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if expectedRequest.IsWebpage {
		if expectedRequest.UseBrowser {
			err = services.DownloadFromDynamic(ctx, expectedRequest.Url, w)
			if err != nil {
				errorHandler(&res, err)
				return
			}
		} else {
			err = services.DownloadFromStatic(ctx, expectedRequest.Url, w)
			if err != nil {
				errorHandler(&res, err)
				return
			}
		}
	} else {
		err = services.DefaultDownloader(ctx, expectedRequest.Url, w)
		if err != nil {
			errorHandler(&res, err)
			return
		}
	}

	if err != nil {
		errorHandler(&res, err)
		return
	}

}

func SetupRoutes(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/download", defaultDownload)

	slog.Info("Listening on 0.0.0.0:8080")
	err := http.ListenAndServe("0.0.0.0:8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server closed")
	} else if err != nil {
		slog.Error(fmt.Sprintf("Unknown error: %s", err))
		os.Exit(1)
	}
}
