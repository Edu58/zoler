package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Edu58/zoler/internal"
)

type CrawlRequest struct {
	URL   string `json:"url"`
	Depth int    `json:"depth"`
}

type CrawlResult struct {
	Took   string `json:"took"`
	Result string `json:"result"`
}

func Crawler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var crawlRequest CrawlRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&crawlRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fetcher := internal.NewFetcher(crawlRequest.URL, crawlRequest.Depth)

	startTime := time.Now()

	result, err := fetcher.Fetch()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timeSince := time.Since(startTime)

	json_response, err := json.Marshal(CrawlResult{
		Took:   timeSince.String(),
		Result: string(result),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(json_response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
