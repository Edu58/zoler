package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/Edu58/zoler/internal/store"
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

func CrawlURL(ctx context.Context, url string) (result []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return []byte{}, fmt.Errorf("error preparing request for URL %s - %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return []byte{}, fmt.Errorf("error getting URL %s - %w", url, err)
	}

	defer func() {
		err = resp.Body.Close()
	}()

	result, err = io.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, fmt.Errorf("error extracting body error: %w", err)
	}

	return result, nil
}

func ProcessResult(store *store.Store, result chan Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for result := range result {
		if result.err != nil {
			log.Printf("Worker %d failed on %s: %v", result.workerId, result.task, result.err)
			continue
		}

		log.Printf("Worker %d crawled %s: %d bytes", result.workerId, result.task, len(result.data))

		if err := store.Set(result.task, string(result.data)); err != nil {
			log.Fatalf("Error setting key: %v", err)
		}

		log.Printf("Saved results for task %s", result.task)
	}
}
