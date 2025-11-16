package handlers

import (
	"net/http"

	"github.com/AnaMaghear/urlshortener/api/internal/models"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

// QR returns a PNG QR code for a short URL.
// Usage: GET http://localhost:8080/qr?code=ViNc94C
func QR(db *gorm.DB, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only GET is allowed
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read code from query: /qr?code=ViNc94C
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "code is required", http.StatusBadRequest)
			return
		}

		// Check that this short URL exists
		var short models.ShortURL
		if err := db.Where("code = ?", code).First(&short).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "short url not found", http.StatusNotFound)
			} else {
				http.Error(w, "database error", http.StatusInternalServerError)
			}
			return
		}

		shortURL := baseURL + "/" + code

		// Generate PNG bytes (256x256)
		png, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
		if err != nil {
			http.Error(w, "failed to generate qr", http.StatusInternalServerError)
			return
		}

		// Send PNG to the client
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(png)
	}
}
