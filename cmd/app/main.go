package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/etsrc/goprod/internal/infra/persistence"
	"github.com/etsrc/goprod/internal/infra/transport/rest"
	"github.com/etsrc/goprod/internal/infra/transport/rest/gen"
	"github.com/etsrc/goprod/internal/service"
)

func main() {
	//lint:ignore SA1019
	bookmarkRepo := persistence.NewInMemoryBookmarkRepository()
	bookmarkService := service.NewBookmarkService(bookmarkRepo)
	handler := rest.NewBookmarkHandler(bookmarkService)

	mux := http.NewServeMux()
	gen.HandlerFromMux(handler, mux)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("🚀 Server starting on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	<-stop
	fmt.Println("\nRestoring peace and quiet... (Shutting down)")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Server exited properly")
}
