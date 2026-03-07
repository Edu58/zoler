package router

import (
	"sync"

	"github.com/Edu58/zoler/internal"
	"github.com/Edu58/zoler/internal/store"
)

var urls = []string{
	"https://www.google.com",
	"https://www.github.com",
	"https://www.stackoverflow.com",
	"https://www.wikipedia.org",
	"https://www.reddit.com",
	"https://www.mozilla.org",
	"https://www.python.org",
	"https://www.postgresql.org",
	"https://www.rust-lang.org",
	"https://www.docker.com",
	"https://www.cloudflare.com",
	"https://www.elastic.co",
	"https://www.phoenixframework.org",
	"https://hex.pm",
	"https://www.erlang.org",
	"https://www.kubernetes.io",
	"https://www.openai.com",
	"https://www.gnu.org",
	"https://www.w3.org",
	"https://www.ietf.org",
}

func Start(db *store.Store, wg *sync.WaitGroup) {
	pool := internal.NewWorkerPool(db, 5)
	go pool.SubmitTasks(urls)
	wg.Add(1)
	go internal.ProcessResult(pool.Store, pool.Results, wg)

}
