package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"solidgate-test-task/card"
)

func main() {
	mux := http.NewServeMux()

	h := InitCardHandler()

	mux.HandleFunc("/validate", h.ValidateCard)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	server := InitServer(mux)

	go func(server *http.Server) {
		log.Printf("server starts listening on: %s", server.Addr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintln(os.Stderr, err)
			stop()
			return
		}
	}(server)

	<-ctx.Done()

	if err := server.Close(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("server closed")
}

func InitCardHandler() *card.Handler {
	return &card.Handler{}
}

func InitServer(mux http.Handler) *http.Server {
	return &http.Server{
		Handler: mux,
		Addr:    "localhost:8080",
	}
}
