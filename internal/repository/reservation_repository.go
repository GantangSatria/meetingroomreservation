package repository

import (
	"meetingroomreservation/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(res *models.Reservation) error
	FindAll() ([]models.Reservation, error)
	FindByID(id uint64) (*models.Reservation, error)
	Update(res *models.Reservation) error
	Delete(id uint64) error
	IsTimeConflict(roomID uint64, start, end time.Time) (bool, error)
	FindByUserRoomAndStartTime(userID uint64, roomID uint64, startTime time.Time) (*models.Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db}
}

func (r *reservationRepository) Create(res *models.Reservation) error {
	return r.db.Create(res).Error
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").Preload("Room").Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByID(id uint64) (*models.Reservation, error) {
	var res models.Reservation
	err := r.db.Preload("User").Preload("Room").First(&res, id).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *reservationRepository) Update(res *models.Reservation) error {
	return r.db.Save(res).Error
}

func (r *reservationRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Reservation{}, id).Error
}

func (r *reservationRepository) IsTimeConflict(roomID uint64, start, end time.Time) (bool, error) {
	var count int64

	err := r.db.Model(&models.Reservation{}).
		Where("room_id = ?", roomID).
		Where("start_time < ? AND end_time > ?", end, start).
		Count(&count).Error

	return count > 0, err
}

func (r *reservationRepository) FindByUserRoomAndStartTime(userID, roomID uint64, startTime time.Time) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.Where("user_id = ? AND room_id = ? AND start_time = ?", userID, roomID, startTime).First(&reservation).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

