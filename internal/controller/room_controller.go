package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"meetingroomreservation/internal/services"
	"meetingroomreservation/pkg/dto"
)

type RoomController struct {
	service services.RoomService
}

func NewRoomController(s services.RoomService) *RoomController {
	return &RoomController{service: s}
}

func (c *RoomController) Create(ctx *gin.Context) {
	var req dto.CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.service.Create(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (c *RoomController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, _ := strconv.ParseUint(idParam, 10, 64)

	var req dto.UpdateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Update(id, &req); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (c *RoomController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, _ := strconv.ParseUint(idParam, 10, 64)

	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (c *RoomController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, _ := strconv.ParseUint(idParam, 10, 64)

	room, err := c.service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	ctx.JSON(http.StatusOK, room)
}

func (c *RoomController) GetAll(ctx *gin.Context) {
	rooms, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, rooms)
}
