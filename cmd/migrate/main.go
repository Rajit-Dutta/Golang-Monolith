package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Rajit-Dutta/GolangMonolith/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	fmt.Println(os.Args)

	cfg := config.MustLoad()

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("migration: NEW %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("migrate: <up |down>")
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown command : %v", err)
	}
}
