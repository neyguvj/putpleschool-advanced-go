package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
)

func randomHandler(w http.ResponseWriter, r *http.Request) {
	num := rand.IntN(6) + 1
	fmt.Fprintf(w, "%d", num)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/random", randomHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server eeror: %v", err)
	}
}
