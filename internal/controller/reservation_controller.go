package controller

import (
	"meetingroomreservation/internal/services"
	"meetingroomreservation/pkg/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReservationController struct {
	svc services.ReservationService
}

func NewReservationController(svc services.ReservationService) *ReservationController {
	return &ReservationController{svc: svc}
}

func (ctrl *ReservationController) Create(c *gin.Context) {
	var req dto.CreateReservationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint64("user_id")

	res, err := ctrl.svc.Create(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (ctrl *ReservationController) GetAll(c *gin.Context) {
	reservations, err := ctrl.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

func (ctrl *ReservationController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	res, err := ctrl.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reservation not found"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (ctrl *ReservationController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var req dto.UpdateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.svc.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reservation updated"})
}

func (ctrl *ReservationController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	if err := ctrl.svc.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reservation deleted"})
}

func (ctrl *ReservationController) GetQRCode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation ID"})
		return
	}

	qrBytes, err := ctrl.svc.GetQRCode(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Writer.Write(qrBytes)
}

func (ctrl *ReservationController) Approve(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can approve reservations"})
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	if err := ctrl.svc.ApproveReservation(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reservation approved"})
}

func (ctrl *ReservationController) Reject(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can reject reservations"})
		return
	}

	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	if err := ctrl.svc.RejectReservation(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reservation rejected"})
}