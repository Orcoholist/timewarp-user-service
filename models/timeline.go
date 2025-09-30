package models

import (
	"gorm.io/gorm"
)

type Timeline struct {
	gorm.Model				   `gorm:"table:directions"`
	ID               uint      `gorm:"primaryKey"`
		Name  string  `gorm:"column:name"`
	Description      string      `gorm:"column:description"`
	Year  int     `gorm:"column:year"`
	// BaseTimelineID   *uint     `gorm:"index" json:"-"`
	// BaseTimeline     *Timeline `json:"-"`
	// CreatedAt        time.Time `json:"created_at"`
	// Events           []Event   `gorm:"foreignKey:TimelineID"`
	
}