package main

import (
	"log"
	"net/http"
	"github.com/tmadelsk/mock-alerts-api/handlers"
)

func main() {
	http.HandleFunc("/alerts", handlers.AlertsHandler)
	log.Println("Mock Alerts API listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
