package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"hnjobs/services/db"
	"hnjobs/services/handler"
	"hnjobs/services/hnjobs"
)

const (
	postgresURI = "postgres://<username>:<password>@<host>:<port>/<database>"
)

func main() {
	migrationVersion, _ := strconv.ParseUint(os.Getenv("MIGRATION_VERSION"), 10, 64)
	log.Printf("Migration version: %d\n", migrationVersion)

	pgURI := os.Getenv("POSTGRES_URI")
	if pgURI == "" {
		pgURI = postgresURI
	}

	db.MigrateDatabase(pgURI, "", uint(migrationVersion))

	// Add HN: Who is hiring IDs
	go hnjobs.FetchStory("20325925") // July 2019
	// go hnjobs.FetchStory("20083795") // June 2019
	// go hnjobs.FetchStory("19797594") // May 2019

	var wg sync.WaitGroup
	wg.Add(1)
	go server(&wg)
	wg.Wait()
}

func server(wg *sync.WaitGroup) {
	http.HandleFunc("/", handler.HandleFunc())
	http.ListenAndServe(":9999", nil)
	wg.Done()
}
