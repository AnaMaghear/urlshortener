package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/AnaMaghear/urlshortener/api/internal/models"
	"gorm.io/gorm"
)

// Usage: POST http://localhost:8080/shorten
type ShortenRequest struct {
	URL       string     `json:"url"`
	Custom    string     `json:"custom"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	Code     string `json:"code"`
}

func Shorten(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Only allow POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		//Parse JSON body
		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		//Validate required fields
		if req.URL == "" {
			http.Error(w, "url is required", http.StatusBadRequest)
			return
		}

		req.URL = strings.TrimSpace(req.URL)

		// generate or use custom code
		code := req.Custom
		if code != "" {
			code = strings.ToLower(code)
		} else {
			code = generateCode()
		}

		// collision check
		var existing models.ShortURL
		result := db.Where("code = ?", code).First(&existing)
		if result.Error == nil {
			http.Error(w, "code already exists", http.StatusConflict)
			return
		}

		var expires *time.Time
		if req.ExpiresAt != nil {
			t := req.ExpiresAt.UTC() // force UTC before saving
			expires = &t
		}

		//Save to database
		newURL := models.ShortURL{
			Code:        code,
			OriginalURL: req.URL,
			ExpiresAt:   expires,
		}

		if err := db.Create(&newURL).Error; err != nil {
			log.Println("DB error:", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		//Prepare response
		resp := ShortenResponse{
			ShortURL: "http://localhost:8080/" + code,
			Code:     code,
		}

		//Return JSON to user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func generateCode() string {
	//This is a string containing all characters we allow in the short code
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 7)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
