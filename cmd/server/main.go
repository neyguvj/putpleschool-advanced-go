package main

import (
	"errors"
	"log"
	"net/http"

	"advancedgo/internal/migrations"
	"advancedgo/internal/product"
	"advancedgo/pkg/db"
	"advancedgo/pkg/logging"
)

func main() {
	cfg, err := db.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gdb, err := db.New(cfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	log.Println("connected to postgres")

	if err := migrations.Migrate(gdb); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	logger := logging.NewLogger()

	productHandler := product.NewHandler(product.NewRepository(gdb))

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	productHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: logging.Middleware(logger)(mux),
	}

	log.Println("listening on :8080")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("serve: %v", err)
	}
}
