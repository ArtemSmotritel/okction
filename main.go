package main

import (
	"flag"
	"fmt"
	"github.com/artemsmotritel/oktion/api"
	"github.com/artemsmotritel/oktion/storage"
	"log"
)

func main() {
	fmt.Println("Hello oktion!")

	var (
		seed    bool
		address string
	)

	flag.BoolVar(&seed, "seed", false, "seed some values in the database")
	flag.StringVar(&address, "address", ":3000", "server address")
	flag.Parse()

	store := storage.NewInMemoryStore()

	if seed {
		log.Println("Seeding data into the database...")
		if err := store.SeedData(); err != nil {
			log.Fatal("Couldn't seed data into the database\n", err.Error())
		}
		log.Println("Finished seeding data into the database")
	}

	server := api.NewServer(address, store)
	log.Println("Listening on", address)
	if err := server.Start(); err != nil {
		log.Fatal(err.Error())
	}
}
