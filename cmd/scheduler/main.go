package main

import (
	"log"
	"os"

	"github.com/Clankyyy/scheduler/internal/metrics"
	"github.com/Clankyyy/scheduler/internal/mgstorage"
	"github.com/Clankyyy/scheduler/internal/routes"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	uri, ok := os.LookupEnv("MONGODB_URI")
	if !ok {
		log.Fatal("Fail to load mongo uri")
	}
	mongoStorage := mgstorage.NewMGStorage(uri)
	prometheusGatherer := metrics.NewHttpRequestsCounter()
	apiServer := routes.NewAPIServer(":8000", mongoStorage, prometheusGatherer) // 8000
	apiServer.Run()
}
