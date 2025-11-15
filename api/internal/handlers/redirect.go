package handlers

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/AnaMaghear/urlshortener/api/internal/geo"
	"github.com/AnaMaghear/urlshortener/api/internal/models"
	"gorm.io/gorm"
)

// Usage: POST http://localhost:8080/ViNc94C
func RedirectHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		code := strings.TrimPrefix(r.URL.Path, "/")
		if code == "" || code == "shorten" || code == "analytics" || code == "qr" {
			http.NotFound(w, r)
			return
		}

		// Look up the original URL in the DB
		var short models.ShortURL
		if err := db.Where("code = ?", code).First(&short).Error; err != nil {
			http.NotFound(w, r)
			return
		}

		now := time.Now().UTC()
		if short.ExpiresAt != nil && now.After(short.ExpiresAt.UTC()) {
			http.Error(w, "link expired", http.StatusGone) // 410 Gone
			return
		}

		// Extract IP (cleaner than r.RemoteAddr)
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr // fallback
		}

		// GEO lookup
		country := geo.LookupCountry(ip)

		// Save click event
		event := models.ClickEvent{
			ShortURLID: short.ID,
			IPAddress:  r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			Country:    country, // later we do geo lookup
			ClickedAt:  time.Now(),
		}

		if err := db.Create(&event).Error; err != nil {
			log.Println("Failed to save click event:", err)
		}

		// Redirect
		http.Redirect(w, r, short.OriginalURL, http.StatusFound)
	}
}
