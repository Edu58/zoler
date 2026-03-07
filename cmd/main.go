package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/Edu58/zoler/internal/store"
	"github.com/Edu58/zoler/router"
)

const (
	STORE_NAME = "./tmp/store"
)

func main() {
	fmt.Println("Starting server")

	var wg sync.WaitGroup
	db, err := store.NewStore(STORE_NAME)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := db.Close()
		log.Fatal(err)
	}()

	router.Start(db, &wg)
	wg.Wait()

	log.Println("Crawling done")
}
