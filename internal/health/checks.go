package health

import (
	"encoding/json"
	"net/http"
	"time"
)

func Handler(start time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "uptime": time.Since(start).String()})
	}
}
