package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
}

func main() {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Configuration error: %s", err)
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Could not initialize database %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Could not close database connection: %v", err)
		}
	}()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		if err := db.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			if encodeErr := json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"detail": "Database connection failed",
			}); encodeErr != nil {
				log.Printf("JSON encoding error: %v", err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if encodeErr := json.NewEncoder(w).Encode(map[string]string{
			"status":   "healthy",
			"database": "connected",
		}); encodeErr != nil {
			log.Printf("JSON encoding error: %v", err)
		}
	})

	log.Println("Server running at port :8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
