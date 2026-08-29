package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Rajit-Dutta/GolangMonolith/internal/config"
	"github.com/Rajit-Dutta/GolangMonolith/internal/db"
	"github.com/Rajit-Dutta/GolangMonolith/internal/handlers"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", handlers.Listings(db))
	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server not running")
	}
}
