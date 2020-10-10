package main

import (
	"flag"
	"log"

	"github.com/muffix/relayr-challenge/internal/database"
	"github.com/muffix/relayr-challenge/internal/httpapi"

	_ "github.com/mattn/go-sqlite3"
)

var (
	databasePath = "offers.db"
	defaultPort  = 8080

	servicePort int
)

func processCommandlineArgs() {
	flag.IntVar(&servicePort, "p", defaultPort, "Port to listen on to serve HTTP requests")
	flag.Parse()
}

func main() {
	processCommandlineArgs()
	service := httpapi.NewService(servicePort)

	db, err := database.InitSQLiteDatabase(databasePath)
	if err != nil {
		log.Fatalf("failed to initialise database: %v", err)
	}

	service.SetDatabase(db)
	service.Start()
}
