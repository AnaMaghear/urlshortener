package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AnaMaghear/urlshortener/api/internal/models"
	"gorm.io/gorm"
)

// Usage:GET http://localhost:8080/analytics?code=ViNc94C
// What we send back to the client.
type AnalyticsResponse struct {
	Code            string           `json:"code"`
	OriginalURL     string           `json:"original_url"`
	TotalClicks     int64            `json:"total_clicks"`
	UniqueIPs       int64            `json:"unique_ips"`
	ClicksByCountry map[string]int64 `json:"clicks_by_country"`
}

func Analytics(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only GET is allowed
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read code from query parameter: /analytics?code=abc123
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "code is required", http.StatusBadRequest)
			return
		}

		// Find the short URL by code
		var short models.ShortURL
		if err := db.Where("code = ?", code).First(&short).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "short url not found", http.StatusNotFound)
			} else {
				http.Error(w, "database error", http.StatusInternalServerError)
			}
			return
		}

		// Total clicks
		var totalClicks int64
		db.Model(&models.ClickEvent{}).
			Where("short_url_id = ?", short.ID).
			Count(&totalClicks)

		// 3Unique IPs
		var uniqueIPs int64
		db.Model(&models.ClickEvent{}).
			Where("short_url_id = ?", short.ID).
			Distinct("ip_address").
			Count(&uniqueIPs)

		// Clicks by country
		type CountryCount struct {
			Country string
			Count   int64
		}
		var rows []CountryCount

		db.Model(&models.ClickEvent{}).
			Select("country, COUNT(*) as count").
			Where("short_url_id = ?", short.ID).
			Group("country").
			Scan(&rows)

		clicksByCountry := map[string]int64{}
		for _, row := range rows {
			country := row.Country
			if country == "" {
				country = "Unknown"
			}
			clicksByCountry[country] = row.Count
		}

		// Build response
		resp := AnalyticsResponse{
			Code:            short.Code,
			OriginalURL:     short.OriginalURL,
			TotalClicks:     totalClicks,
			UniqueIPs:       uniqueIPs,
			ClicksByCountry: clicksByCountry,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
