// handlers/directions.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"supabase-go-server/models"

	"gorm.io/gorm"
)

func GetDirections(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var directions []models.Timeline
		if err := db.Raw("SELECT * FROM directions").Scan(&directions).Error; err != nil {
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			log.Printf("Query failed: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(directions); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Encoding failed: %v", err)
		}
	}
}
