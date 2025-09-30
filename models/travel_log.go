// models/travel_log.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type TravelLog struct {
	gorm.Model `gorm:"table:travel_logs"`
	ID           uint      `gorm:"primaryKey"`
	TimeTravelID uint      `gorm:"index" json:"-"`
	TimeTravel   TimeTravel `json:"-"`
	EventTime    time.Time  `json:"event_time"`
	Details      string     `json:"details"`
	Location     string     `json:"location"`
}
