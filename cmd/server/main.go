package main

import (
	"log"

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
	log.Println("connected to postgres")

	var count int64
	if err := gdb.Table("products").Count(&count).Error; err != nil {
		log.Fatalf("select from products: %v", err)
	}
	log.Printf("select ok, products: %d", count)
}
