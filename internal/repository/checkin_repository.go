package repository

import (
	"meetingroomreservation/internal/models"
	"gorm.io/gorm"
)

type CheckinRepository interface {
	Create(checkin *models.Checkin) error
	FindByReservationID(reservationID uint64) (*models.Checkin, error)
	Update(checkin *models.Checkin) error
}

type checkinRepository struct {
	db *gorm.DB
}

func NewCheckinRepository(db *gorm.DB) CheckinRepository {
	return &checkinRepository{db}
}

func (r *checkinRepository) Create(checkin *models.Checkin) error {
	return r.db.Create(checkin).Error
}

func (r *checkinRepository) FindByReservationID(reservationID uint64) (*models.Checkin, error) {
	var c models.Checkin
	err := r.db.Preload("Reservation").First(&c, "reservation_id = ?", reservationID).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *checkinRepository) Update(checkin *models.Checkin) error {
	return r.db.Save(checkin).Error
}

func (r *checkinRepository) GetByReservationID(reservationID uint64) (*models.Checkin, error) {
	var checkin models.Checkin
	if err := r.db.Where("reservation_id = ?", reservationID).First(&checkin).Error; err != nil {
		return nil, err
	}
	return &checkin, nil
}
