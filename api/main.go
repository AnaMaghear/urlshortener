package main

import (
	"log"
	"net/http"
	"time"

	"github.com/AnaMaghear/urlshortener/api/internal/config"
	"github.com/AnaMaghear/urlshortener/api/internal/handlers"
	"github.com/AnaMaghear/urlshortener/api/internal/middleware"
	"github.com/AnaMaghear/urlshortener/api/internal/models"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	_ = godotenv.Load(".env")

	// Read config
	cfg := config.Load()

	if cfg.DBDSN == "" {
		log.Fatal("DB_DSN is empty in .env")
	}

	// Connect to Postgres using GORM
	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Auto-create tables
	if err := db.AutoMigrate(&models.ShortURL{}, &models.ClickEvent{}); err != nil {
		log.Fatalf("Failed to migrate DB: %v", err)
	}

	log.Println("DB connected and tables created!")

	rl := middleware.NewRateLimiter(10, time.Minute)
	mux := http.NewServeMux()
	rateLimitedHandler := rl.Middleware(mux)

	mux.HandleFunc("/shorten", handlers.Shorten(db))
	mux.HandleFunc("/analytics", handlers.Analytics(db))
	mux.HandleFunc("/qr", handlers.QR(db, cfg.BaseURL))
	mux.HandleFunc("/", handlers.RedirectHandler(db))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Printf("Server running on port %s", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, rateLimitedHandler)
}
