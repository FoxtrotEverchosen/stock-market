package main

import (
	"log"
	"net/http"
	"os"

	rdb "github.com/bulanda/stock-market/src/redis"
	"github.com/bulanda/stock-market/src/routes"
)

func main() {
	rdb.Connect()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	routes.Register(mux)

	log.Printf("[Server] Starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("[Server] Failed to start: %v", err)
	}
}
