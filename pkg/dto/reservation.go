package dto

import "time"

type CreateReservationRequest struct {
	RoomID    uint64    `json:"room_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

type UpdateReservationRequest struct {
	RoomID    uint64    `json:"room_id" binding:"required"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

type UserReservationResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RoomReservationResponse struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	Capacity   int    `json:"capacity"`
	Facilities string `json:"facilities"`
}

type ReservationResponse struct {
	ID        uint64                  `json:"id"`
	User      UserReservationResponse `json:"user"`
	Room      RoomReservationResponse `json:"room"`
	StartTime time.Time               `json:"start_time"`
	EndTime   time.Time               `json:"end_time"`
	Status    string                  `json:"status"`
	QRCode    *string                 `json:"qrcode"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
}