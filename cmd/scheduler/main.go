package main

import (
	"log"
	"os"

	"github.com/Clankyyy/scheduler/internal/mgstorage"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	uri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		log.Fatal("Fail to load mongo uri")
	}
	mongo := mgstorage.NewMGStorage(uri)
	log.Print(mongo)
	// s := storage.NewFSStorage("data/spbgti/")
	// apiServer := routes.NewAPIServer(":8000", s)
	// apiServer.Run()
}
