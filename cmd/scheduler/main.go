package main

import (
	"github.com/Clankyyy/scheduler/internal/routes"
)

func main() {
	apiServer := routes.NewAPIServer(":8000")
	apiServer.Run()
}
