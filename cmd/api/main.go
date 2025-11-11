package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"santori/linkchecker/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.Check(w, r)
		default:
			http.Error(w, "use POST method to send links", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.GenerateReport(w, r)
		default:
			http.Error(w, "use POST method to send links", http.StatusMethodNotAllowed)
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Run server in a goroutine so we can wait for shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
