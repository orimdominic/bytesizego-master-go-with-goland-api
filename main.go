package main

import (
	"first-api/internal/db"
	"first-api/internal/todo"
	"first-api/internal/transport"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("unable to load env variables")
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("invalid env value DB_PORT", err)
	}

	db, err := db.New(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		dbPort,
	)
	if err != nil {
		log.Fatalf("could not connect to database %v", err)
	}

	todoSvc := todo.NewService(db)
	server := transport.NewServer(todoSvc)

	log.Fatal(server.Serve())
}
