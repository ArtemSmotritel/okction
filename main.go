package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/artemsmotritel/oktion/api"
	"github.com/artemsmotritel/oktion/storage"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"log"
)

func main() {
	fmt.Println("Hello oktion!")

	var (
		seed    bool
		address string
		dbURL   string
	)

	flag.BoolVar(&seed, "seed", false, "seed some values in the database")
	flag.StringVar(&address, "address", ":3000", "server address")
	flag.StringVar(&dbURL, "db", "postgres://postgres:abobus@localhost:5432/oktion", "database connection url")
	flag.Parse()

	logger := log.Default()

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	pgxdecimal.Register(conn.TypeMap())

	//store := storage.NewInMemoryStore()
	store := storage.NewPostgresqlStore(conn, logger)

	if seed {
		logger.Println("Seeding data into the database...")
		if err := store.SeedData(); err != nil {
			logger.Fatal("Couldn't seed data into the database\n", err.Error())
		}
		logger.Println("Finished seeding data into the database")
	}

	server := api.NewServer(address, store, logger)
	logger.Println("Listening on", address)
	if err := server.Start(); err != nil {
		logger.Fatal(err.Error())
	}
}
