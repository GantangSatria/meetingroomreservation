package services

import (
	"errors"
	"meetingroomreservation/internal/models"
	"meetingroomreservation/internal/repository"
	"meetingroomreservation/pkg/dto"
	"meetingroomreservation/internal/utils"
)

type UserService interface {
	Register(req *dto.RegisterRequest) (uint64, error)
	Login(req *dto.LoginRequest) (string, error)
	GetAll() ([]models.User, error)
	GetByID(id uint) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
	jwtSecret string
}

func NewUserService(r repository.UserRepository, jwtSecret string) UserService {
	return &userService{repo: r, jwtSecret: jwtSecret}
}

func (s *userService) Register(req *dto.RegisterRequest) (uint64, error) {
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return 0, errors.New("email already registered")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return 0, err
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
	}

	if err := s.repo.Create(user); err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (s *userService) Login(req *dto.LoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}
