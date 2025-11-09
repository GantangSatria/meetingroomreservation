package services

import (
	"errors"
	"meetingroomreservation/internal/models"
	"meetingroomreservation/internal/repository"
	"meetingroomreservation/pkg/dto"
)

type RoomService interface {
	Create(req *dto.CreateRoomRequest) (uint64, error)
	Update(id uint64, req *dto.UpdateRoomRequest) error
	Delete(id uint64) error
	GetByID(id uint64) (*models.Room, error)
	GetAll() ([]models.Room, error)
}

type roomService struct {
	repo repository.RoomRepository
}

func NewRoomService(r repository.RoomRepository) RoomService {
	return &roomService{repo: r}
}

func (s *roomService) Create(req *dto.CreateRoomRequest) (uint64, error) {
	room := &models.Room{
		Name:       req.Name,
		Location:   req.Location,
		Capacity:   req.Capacity,
		Facilities: req.Facilities,
	}
	if err := s.repo.Create(room); err != nil {
		return 0, err
	}
	return room.ID, nil
}

func (s *roomService) Update(id uint64, req *dto.UpdateRoomRequest) error {
	room, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("room not found")
	}

	if req.Name != "" {
		room.Name = req.Name
	}
	if req.Location != "" {
		room.Location = req.Location
	}
	if req.Capacity != 0 {
		room.Capacity = req.Capacity
	}
	if req.Facilities != "" {
		room.Facilities = req.Facilities
	}

	return s.repo.Update(room)
}

func (s *roomService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

func (s *roomService) GetByID(id uint64) (*models.Room, error) {
	return s.repo.FindByID(id)
}

func (s *roomService) GetAll() ([]models.Room, error) {
	return s.repo.FindAll()
}
