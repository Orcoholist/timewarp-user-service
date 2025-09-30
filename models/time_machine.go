package models

import (
	// "time"

	"gorm.io/gorm"
)

type TimeMachine struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	MaxDistance   int       `json:"max_distance"`
	EnergyLevel   int       `json:"energy_level"`
	Status        string    `json:"status"` // active, maintenance
	TimeTravels   []TimeTravel `gorm:"foreignKey:TimeMachineID"`
}