package main

import (
	"go-order-inventory/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
