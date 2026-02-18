package internal

import (
	"fmt"
	"io"
	"net/http"
)

type Fetcher struct {
	url   string
	depth int
}

func NewFetcher(url string, depth int) *Fetcher {
	return &Fetcher{url, depth}
}

func (f *Fetcher) Fetch() ([]byte, error) {
	resp, err := http.Get(f.url)

	if err != nil {
		return []byte{}, fmt.Errorf("error calling URL %s - %w", f.url, err)
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, fmt.Errorf("error extracting body error: %w", err)
	}

	return result, nil
}

func CrawlURL(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return []byte{}, fmt.Errorf("error calling URL %s - %w", url, err)
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, fmt.Errorf("error extracting body error: %w", err)
	}

	return result, nil
}
