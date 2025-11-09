package models

import "time"

type Room struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	Name       string    `gorm:"size:100;not null"`
	Location   string    `gorm:"size:100;not null"`
	Capacity   int       `gorm:"not null"`
	Facilities string    `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
