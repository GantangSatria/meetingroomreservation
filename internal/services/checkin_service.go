package services

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"meetingroomreservation/internal/models"
	"meetingroomreservation/internal/repository"
)

type CheckinService interface {
	Checkin(reservationID uint64) (*models.Checkin, error)
	Checkout(reservationID uint64) (*models.Checkin, error)
	CheckinByQRCode(qrData string) (*models.Checkin, error)
}

type checkinService struct {
	repoCheckin     repository.CheckinRepository
	repoReservation repository.ReservationRepository
}

func NewCheckinService(checkinRepo repository.CheckinRepository, reservationRepo repository.ReservationRepository) CheckinService {
	return &checkinService{
		repoCheckin:     checkinRepo,
		repoReservation: reservationRepo,
	}
}

func (s *checkinService) Checkin(reservationID uint64) (*models.Checkin, error) {
	res, err := s.repoReservation.FindByID(reservationID)
	if err != nil {
		return nil, errors.New("reservation not found")
	}

	if res.Status != "approved" {
		return nil, errors.New("cannot check in, reservation not approved")
	}

	existing, _ := s.repoCheckin.FindByReservationID(reservationID)
	if existing != nil {
		return nil, errors.New("already checked in")
	}

	checkin := &models.Checkin{
		ReservationID: reservationID,
		CheckinTime:   time.Now(),
	}

	if err := s.repoCheckin.Create(checkin); err != nil {
		return nil, err
	}

	res.Status = "checked_in"
	_ = s.repoReservation.Update(res)

	return checkin, nil
}

func (s *checkinService) Checkout(reservationID uint64) (*models.Checkin, error) {
	checkin, err := s.repoCheckin.FindByReservationID(reservationID)
	if err != nil {
		return nil, errors.New("checkin not found")
	}

	if checkin.CheckoutTime != nil {
		return nil, errors.New("already checked out")
	}

	now := time.Now()
	checkin.CheckoutTime = &now

	if err := s.repoCheckin.Update(checkin); err != nil {
		return nil, err
	}

	return checkin, nil
}

func (s *checkinService) CheckinByQRCode(qrData string) (*models.Checkin, error) {
	// masih error :(
	parts := strings.Split(qrData, ":")
	if len(parts) < 4 || parts[0] != "reservation" {
		return nil, errors.New("invalid QR format")
	}

	userID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, errors.New("invalid user ID in QR")
	}

	roomID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return nil, errors.New("invalid room ID in QR")
	}

	startTimeStr := strings.Join(parts[3:], ":")
	startTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", startTimeStr)
	if err != nil {
		return nil, errors.New("invalid start time format")
	}

	reservation, err := s.repoReservation.FindByUserRoomAndStartTime(userID, roomID, startTime)
	if err != nil {
		return nil, errors.New("reservation not found for this QR")
	}

	if reservation.Status != "approved" {
		return nil, errors.New("reservation not approved yet")
	}

	existing, _ := s.repoCheckin.FindByReservationID(reservation.ID)
	if existing != nil {
		return nil, errors.New("already checked in")
	}

	checkin := &models.Checkin{
		ReservationID: reservation.ID,
		CheckinTime:   time.Now(),
	}

	if err := s.repoCheckin.Create(checkin); err != nil {
		return nil, err
	}

	reservation.Status = "checked_in"
	_ = s.repoReservation.Update(reservation)

	return checkin, nil
}

