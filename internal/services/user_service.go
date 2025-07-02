package services

import (
	"test-task/internal/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByID(guid string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", guid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
