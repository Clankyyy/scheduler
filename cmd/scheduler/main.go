package main

import (
	"github.com/Clankyyy/scheduler/internal/routes"
	"github.com/Clankyyy/scheduler/internal/schedule"
	"github.com/Clankyyy/scheduler/internal/storage"
)

func main() {
	s := storage.NewFSStorage("data/spbgti/")
	apiServer := routes.NewAPIServer(":8000", s)
	apiServer.Run()
}

func testing() {
	schedule.Test()
}
