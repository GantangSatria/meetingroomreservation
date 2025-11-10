package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"meetingroomreservation/internal/models"
	"meetingroomreservation/internal/repository"
	"meetingroomreservation/internal/utils"
	"meetingroomreservation/pkg/dto"
)

type ReservationService interface {
	Create(req *dto.CreateReservationRequest, userID uint64) (*dto.ReservationResponse, error)
	GetByID(id uint64) (*dto.ReservationResponse, error)
	GetAll() ([]dto.ReservationResponse, error)
	Update(id uint64, req *dto.UpdateReservationRequest) error
	Delete(id uint64) error
	GetQRCode(id uint64) ([]byte, error)

	ApproveReservation(id uint64) error
	RejectReservation(id uint64) error
}

type reservationService struct {
	repo repository.ReservationRepository
}

func NewReservationService(r repository.ReservationRepository) ReservationService {
	return &reservationService{repo: r}
}

func (s *reservationService) Create(req *dto.CreateReservationRequest, userID uint64) (*dto.ReservationResponse, error) {
	conflict, err := s.repo.IsTimeConflict(req.RoomID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}

	if conflict {
		return nil, errors.New("time slot is already booked")
	}

	res := &models.Reservation{
		UserID:    userID,
		RoomID:    req.RoomID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Status:    "pending",
	}

	qrData := fmt.Sprintf("reservation:%d:%d:%s", userID, req.RoomID, req.StartTime.String())
	qr, err := utils.GenerateQRCode(qrData)
	if err != nil {
		return nil, err
	}
	res.QRCode = &qr

	if err := s.repo.Create(res); err != nil {
		return nil, err
	}

	res, err = s.repo.FindByID(res.ID)
	if err != nil {
		return nil, err
	}

	return toReservationResponse(res), nil
}

func (s *reservationService) GetByID(id uint64) (*dto.ReservationResponse, error) {
	res, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return toReservationResponse(res), nil
}


func (s *reservationService) GetAll() ([]dto.ReservationResponse, error) {
	reservations, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []dto.ReservationResponse
	for _, res := range reservations {
		responses = append(responses, *toReservationResponse(&res))
	}

	return responses, nil
}

func (s *reservationService) Update(id uint64, req *dto.UpdateReservationRequest) error {
	res, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	conflict, err := s.repo.IsTimeConflict(req.RoomID, req.StartTime, req.EndTime)
	if err != nil {
		return err
	}
	if conflict {
		return errors.New("time slot is already booked")
	}

	res.RoomID = req.RoomID
	res.StartTime = req.StartTime
	res.EndTime = req.EndTime
	res.Status = req.Status

	return s.repo.Update(res)
}

func (s *reservationService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

func (s *reservationService) GetQRCode(id uint64) ([]byte, error) {
	res, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if res.QRCode == nil {
		return nil, errors.New("QR code not found for this reservation")
	}

	qrBytes, err := base64.StdEncoding.DecodeString(*res.QRCode)
	if err != nil {
		return nil, errors.New("failed to decode QR code")
	}

	return qrBytes, nil
}

func toReservationResponse(res *models.Reservation) *dto.ReservationResponse {
	return &dto.ReservationResponse{
		ID: res.ID,
		User: dto.UserReservationResponse{
			ID:    res.User.ID,
			Name:  res.User.Name,
			Email: res.User.Email,
		},
		Room: dto.RoomReservationResponse{
			ID:         res.Room.ID,
			Name:       res.Room.Name,
			Location:   res.Room.Location,
			Capacity:   res.Room.Capacity,
			Facilities: res.Room.Facilities,
		},
		StartTime: res.StartTime,
		EndTime:   res.EndTime,
		Status:    res.Status,
		QRCode:    res.QRCode,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}
}

func (s *reservationService) ApproveReservation(id uint64) error {
	res, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if res.Status != "pending" {
		return errors.New("only pending reservations can be approved")
	}

	res.Status = "approved"
	return s.repo.Update(res)
}

func (s *reservationService) RejectReservation(id uint64) error {
	res, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if res.Status != "pending" {
		return errors.New("only pending reservations can be rejected")
	}

	res.Status = "rejected"
	return s.repo.Update(res)
}
