package geo

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// LookupCountry returns a 2-letter country code (e.g. "US", "RO").
// If lookup fails, returns "--".
func LookupCountry(ip string) string {
	client := http.Client{
		Timeout: 1 * time.Second, // avoid blocking on slow API
	}

	// Call ipapi.co
	resp, err := client.Get("https://ipapi.co/" + ip + "/country/")
	if err != nil {
		return "--"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "--"
	}

	country := strings.TrimSpace(string(body))
	if len(country) != 2 {
		return "--"
	}

	return country
}
