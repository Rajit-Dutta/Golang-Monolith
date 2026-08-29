package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Rajit-Dutta/GolangMonolith/internal/config"
	"github.com/Rajit-Dutta/GolangMonolith/internal/db"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.MustLoad()
	_, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("sql.Open: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
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
