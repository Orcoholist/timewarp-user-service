package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"supabase-go-server/models"

	"gorm.io/gorm"
)

func GetUserHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			log.Printf("Query failed: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Encoding failed: %v", err)
		}
	}
}
