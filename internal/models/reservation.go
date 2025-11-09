package models

import "time"

type Reservation struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `gorm:"not null"`
	RoomID    uint64    `gorm:"not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Status    string    `gorm:"type:varchar(20);default:'pending';not null"`
	QRCode    *string   `gorm:"type:text"`

	User User `gorm:"foreignKey:UserID"`
	Room Room `gorm:"foreignKey:RoomID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
