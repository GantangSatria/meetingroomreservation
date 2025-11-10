package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"meetingroomreservation/internal/services"
)

type CheckinController struct {
	svc services.CheckinService
}

func NewCheckinController(svc services.CheckinService) *CheckinController {
	return &CheckinController{svc: svc}
}

func (ctrl *CheckinController) Checkin(c *gin.Context) {
	idStr := c.Param("reservation_id")
	reservationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	checkin, err := ctrl.svc.Checkin(reservationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Check-in successful",
		"data":    checkin,
	})
}

func (ctrl *CheckinController) Checkout(c *gin.Context) {
	idStr := c.Param("reservation_id")
	reservationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	checkin, err := ctrl.svc.Checkout(reservationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Checkout successful",
		"data":    checkin,
	})
}

func (ctrl *CheckinController) CheckinByQRCode(c *gin.Context) {
	var req struct {
		QRData string `json:"qr_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	checkin, err := ctrl.svc.CheckinByQRCode(req.QRData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Check-in successful via QR code",
		"data":    checkin,
	})
}
