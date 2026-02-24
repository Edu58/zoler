package internal

import (
	"context"
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

type Result struct {
	workerId int
	task     string
	data     []byte
	err      error
}

type WorkerPool struct {
	workers []*Worker
	Results chan Result
	ctx     context.Context
	cancel  context.CancelFunc
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

func (wp *Worker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Worker %d shitting down: %v", wp.id, ctx.Err())
				return
			default:
			}

			if task, err := wp.queue.dequeue(); err == nil {
				process_task(ctx, task, wp.id, wp.pool.Results)
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()
}

func process_task(ctx context.Context, task string, workerId int, resultChan chan<- Result) {
	taskCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	result, err := CrawlURL(taskCtx, task)
	cancel()

	if err != nil {
		log.Printf("Worker: %d task failed for url %s with error: %v", workerId, task, err)
	}

	select {
	case resultChan <- Result{workerId: workerId, task: task, err: err, data: result}:
	case <-ctx.Done(): // don't block if the pool is shutting down
		return
	}
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

// Shutdown gracefully stops all workers
func (wp *WorkerPool) Shutdown() {
	wp.cancel()
}

func NewWorkerPool(numOfWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{ctx: ctx,
		cancel:  cancel,
		workers: make([]*Worker, numOfWorkers),
		Results: make(chan Result, 1000),
	}

	for i := range numOfWorkers {
		pool.workers[i] = &Worker{id: i, queue: &URLQueue{}, pool: pool}
		pool.workers[i].Start(ctx)
	}

	return pool
}
