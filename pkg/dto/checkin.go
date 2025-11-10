package dto

import "time"

type CheckinResponse struct {
	ID            uint64                 `json:"id"`
	ReservationID uint64                 `json:"reservation_id"`
	CheckinTime   time.Time              `json:"checkin_time"`
	CheckoutTime  *time.Time             `json:"checkout_time,omitempty"`
	Reservation   *ReservationSimpleInfo `json:"reservation,omitempty"`
}

type ReservationSimpleInfo struct {
	ID        uint64    `json:"id"`
	RoomName  string    `json:"room_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}