package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"meetingroomreservation/pkg/dto"
	"meetingroomreservation/internal/services"
)

type UserController struct {
	svc services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{svc: s}
}

func (uc *UserController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uc.svc.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (uc *UserController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := uc.svc.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{Token: token})
}

func (uc *UserController) GetAll(c *gin.Context) {
	users, err := uc.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		out = append(out, map[string]interface{}{
			"id": u.ID, "name": u.Name, "email": u.Email, "created_at": u.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func GetUserID(c *gin.Context) (uint, bool) {
	uidVal, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	switch v := uidVal.(type) {
	case uint:
		return v, true
	case int:
		return uint(v), true
	case float64:
		return uint(v), true
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return uint(i), true
		}
	}
	return 0, false
}
