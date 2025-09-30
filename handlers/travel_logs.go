// handlers/travel_logs.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"supabase-go-server/models"

	"gorm.io/gorm"
)

func GetTravelLogs(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var travelLogs []models.TravelLog
		if err := db.Debug().Find(&travelLogs).Error; err != nil {
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			log.Printf("Query failed: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(travelLogs); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Encoding failed: %v", err)
		}
	}
}
