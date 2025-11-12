package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"santori/linkchecker/handlers"
	"santori/linkchecker/models"
	"syscall"
	"time"
)

func main() {

	models.StartWorkerPool(5)

	mux := http.NewServeMux()
	mux.HandleFunc("/check", handlers.Check)
	mux.HandleFunc("/report", handlers.GenerateReport)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	models.StopWorkerPool()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server is off due to: %v", err)
	}

	log.Println("Server is off :)")
}
