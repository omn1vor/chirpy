package main

import (
	"flag"
	"log"
	"os"

	"github.com/omn1vor/chirpy/internal/database"
)

const dbPath = "database.json"

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	serviceId      string
	jwtSecret      string
}

func newApiConfig() *apiConfig {
	truncateDbIfNeeded()

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal("Error creating a database file:", err.Error())
	}

	cfg := apiConfig{
		db:        db,
		serviceId: "chirpy",
		jwtSecret: os.Getenv("JWT_SECRET"),
	}
	return &cfg
}

func truncateDbIfNeeded() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debug {
		os.Remove(dbPath)
	}
}
