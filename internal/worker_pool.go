package internal

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

type URLQueue struct {
	mu   sync.Mutex
	URLS []string
}

type WorkerPool struct {
	workers []*Worker
}

type Worker struct {
	id    int
	queue *URLQueue
	pool  *WorkerPool
}

func (uq *URLQueue) enqueue(url string) {
	uq.mu.Lock()
	defer uq.mu.Unlock()

	uq.URLS = append(uq.URLS, url)
}

func (uq *URLQueue) dequeue() (string, error) {
	uq.mu.Lock()
	defer uq.mu.Unlock()

	if len(uq.URLS) < 1 {
		return "", errors.New("0 urls remaining in queue")
	}

	url := uq.URLS[0]

	uq.URLS = uq.URLS[1:]

	return url, nil
}

func (wp *Worker) Start() {
	go func() {
		for {
			if task, err := wp.queue.dequeue(); err == nil {
				result, _ := CrawlURL(task)
				log.Printf("Worker: %d got: %v", wp.id, string(result))
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()
}

func (wp *WorkerPool) SubmitTask(url string) {
	randomWorker := rand.Intn(len(wp.workers))
	wp.workers[randomWorker].queue.enqueue(url)
}

func (wp *WorkerPool) SubmitTasks(urls []string) {
	for _, url := range urls {
		wp.SubmitTask(url)
	}
}

func NewWorkerPool(numOfWorkers int) *WorkerPool {
	pool := &WorkerPool{workers: make([]*Worker, numOfWorkers)}

	for i := range numOfWorkers {
		pool.workers[i] = &Worker{id: i, queue: &URLQueue{}, pool: pool}
		pool.workers[i].Start()
	}

	return pool
}
