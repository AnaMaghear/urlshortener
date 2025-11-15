package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env
	_ = godotenv.Load(".env")

	// 2. Read env variables
	port := getEnv("PORT", "8080")
	dsn := getEnv("DB_DSN", "")

	if dsn == "" {
		log.Fatal("DB_DSN is empty")
	}

	// 3. Connect to Postgres
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	log.Println("Connected to Postgres successfully")

	// 4. Simple health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	log.Printf("Server listening on :%s ...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
