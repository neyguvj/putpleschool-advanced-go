package main

import (
	"log"

	"advancedgo/internal/migrations"
	"advancedgo/pkg/db"
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

	if err := migrations.Migrate(gdb); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	log.Println("migrations applied")
}
