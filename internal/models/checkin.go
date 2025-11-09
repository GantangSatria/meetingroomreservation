package models

import "time"

type Checkin struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement"`
	ReservationID uint64     `gorm:"not null"`
	CheckinTime   time.Time  `gorm:"not null"`
	CheckoutTime  *time.Time `gorm:"null"`

	Reservation Reservation `gorm:"foreignKey:ReservationID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
