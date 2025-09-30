package models

import (
	"time"

	"gorm.io/gorm"
)

type TimeTravel struct {
	gorm.Model
	UserID           uint      `gorm:"index" json:"-"`
	User             *User     `json:"-"`
	TimeMachineID    uint      `gorm:"index" json:"-"`
	TimeMachine      *TimeMachine `json:"-"`
	StartTime        time.Time `json:"start_time"`
	Destination      time.Time `json:"destination"`
	Origin           string    `json:"origin"`
	DestinationPoint string    `json:"destination_point"`
	Status           string    `json:"status"` // planned, completed, canceled
	Description      string    `json:"description"`
	TimelineID       uint      `gorm:"index" json:"-"`
	Timeline         *Timeline `json:"-"`
	TravelLogs       []TravelLog `gorm:"foreignKey:TimeTravelID"`
}
