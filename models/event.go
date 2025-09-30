package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey"`
	TimelineID  uint      `gorm:"index" json:"-"`
	Timeline    *Timeline `json:"-"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	Impact      string    `json:"impact"`
}
