package main

import (
	"github.com/Clankyyy/scheduler/internal/schedule"
)

func main() {
	testing()
	// apiServer := routes.NewAPIServer(":8000")
	// apiServer.Run()
}

func testing() {
	schedule.Test()
}
