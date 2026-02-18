package router

import (
	"time"

	"github.com/Edu58/zoler/internal"
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

func Start() {
	pool := internal.NewWorkerPool(5)
	pool.SubmitTasks(urls)

	time.Sleep(time.Minute)
	// mux := http.NewServeMux()

	// mux.HandleFunc("/crawl", controller.Crawler)

	// log.Fatal(http.ListenAndServe(":4500", mux))
}
