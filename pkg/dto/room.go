package dto

type CreateRoomRequest struct {
	Name       string `json:"name" binding:"required"`
	Location   string `json:"location" binding:"required"`
	Capacity   int    `json:"capacity" binding:"required,gt=0"`
	Facilities string `json:"facilities"`
}

type UpdateRoomRequest struct {
	Name       string `json:"name"`
	Location   string `json:"location"`
	Capacity   int    `json:"capacity"`
	Facilities string `json:"facilities"`
}
