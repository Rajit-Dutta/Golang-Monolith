package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Rajit-Dutta/GolangMonolith/internal/config"
	"github.com/Rajit-Dutta/GolangMonolith/internal/db"
	"github.com/Rajit-Dutta/GolangMonolith/internal/handlers"
	"github.com/Rajit-Dutta/GolangMonolith/internal/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	loggerHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}

	mux := http.NewServeMux()
	lh := handlers.ListingConstructor(db, logger)

	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	handler := middleware.RequestID(mux)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server not running: %v", err)
	}
}
