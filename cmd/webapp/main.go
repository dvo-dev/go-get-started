package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dvo-dev/go-get-started/pkg/datastorage"
	"github.com/dvo-dev/go-get-started/pkg/handlers"
	"github.com/dvo-dev/go-get-started/pkg/server"
)

func main() {
	log.Println("Starting main webapp server...")

	// This infinite loop is to constantly attempt restarts should run error
	for {
		if err := run(); err != nil {
			log.Printf("Error occurred serving webapp: %v\n", err)
		}
	}
}

// While we could directly init + serve directly in `main()`, having a separate
// function will allow us to call it and reattempt serving should we error out.
//
// This is a simple "dumb" approach for `main() to have multiple attempts at
// initializing, as well as soft restarting the server if it crashes at any
// point.
func run() error {
	var err error
	s := server.Server{}.InitializeServer()

	// TODO: env vars
	webappPort := os.Getenv("WEBAPP_HOST_PORT")
	if len(webappPort) == 0 {
		webappPort = "8080"
	}

	// TODO: add route handlers
	s.AssignHandler(
		"/health",
		handlers.RecoveryWrapper(handlers.HandleHealth()),
	)

	mem := datastorage.MemStorage{}.Initialize()
	dsHandler := handlers.DataStorageHandler{}.Initialize(
		mem,
	)
	s.AssignHandler(
		"/datastorage",
		handlers.RecoveryWrapper(
			dsHandler.HandleClientRequest(),
		),
	)

	log.Println("Webapp server has been initialized, now serving...")
	err = s.ServeAndListen(fmt.Sprintf(":%s", webappPort))

	return err
}
