package main

import (
	"fmt"
	"log"

	"github.com/Edu58/zoler/internal/store"
	"github.com/Edu58/zoler/router"
)

const (
	STORE_NAME = "./tmp/store"
)

func main() {
	fmt.Println("Starting server")

	db, err := store.NewStore(STORE_NAME)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := db.Close()
		log.Fatal(err)
	}()

	if err := db.Set("hello", "world"); err != nil {
		log.Fatalf("Error setting key: %v", err)
	}

	exists, _ := db.Exists("hello")

	fmt.Println("Exists check: ", exists)

	result, err := db.Get("hello")

	if err != nil {
		log.Fatalf("Error setting key: %v", err)
	}

	fmt.Println("Got result: ", result)

	router.Start()
}
