package repository

import (
	"meetingroomreservation/internal/models"

	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *models.Room) error
	Update(room *models.Room) error
	Delete(id uint64) error
	FindByID(id uint64) (*models.Room, error)
	FindAll() ([]models.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

func (r *roomRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Room{}, id).Error
}

func (r *roomRepository) FindByID(id uint64) (*models.Room, error) {
	var room models.Room
	if err := r.db.First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindAll() ([]models.Room, error) {
	var rooms []models.Room
	if err := r.db.Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}
