package main

import (
	"advancedgo/internal/email/config"
	"advancedgo/internal/email/send"
	"advancedgo/internal/email/storage"
	"advancedgo/internal/email/verify"
	"errors"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	storage := storage.NewStorage("storage.json")
	router := http.NewServeMux()
	send.NewSendHandler(router, cfg, storage)
	verify.NewVerifyHandler(router, storage)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
