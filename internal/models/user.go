package models

import "time"

type User struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:100;uniqueIndex;not null"`
	Password 	 string    `gorm:"type:text;not null"`
	Role 		 string    `gorm:"type:text;check:role IN ('user','admin');default:'user'"` 
	CreatedAt    time.Time
	UpdatedAt    time.Time
}