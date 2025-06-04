/* Main (Entrypoint)

- Connects to the mongo db database
- Sets up the routes in a mux
- Starts the server in a goroutine
- Gracefully shuts down the server on termination signal

Rodrigo Núñez, Joaquin Badillo
2025-06-04
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ccs/logger/routes"
	"ccs/logger/db"
)

func Colorize(color int, message string) string {
	return fmt.Sprintf("\033[0;%dm%s\033[0m", color, message)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /logs", routes.PostLog)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	_, err := db.GetMongoClient()
	if err != nil {
		log.Fatalf("failed to connect to MongoDB instance: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("⚡ Server running")
		if os.Getenv("PRODUCTION") == "" {
			log.Println(Colorize(34, fmt.Sprintf("   http://localhost%s", port)))
		}
		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-sigChan
	log.Println("🕊️ Gracefully shutting down")
	db.CloseMongoClient(context.Background())
	server.Shutdown(context.Background())
}
