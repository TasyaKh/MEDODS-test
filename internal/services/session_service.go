package services

import (
	"test-task/internal/config"
	"test-task/internal/dto"
	"test-task/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SessionService struct {
	config *config.Config
	db     *gorm.DB
}

func NewSessionService(db *gorm.DB, config *config.Config) *SessionService {
	return &SessionService{db: db, config: config}
}

func (s *SessionService) AddSession(data dto.SessionCreate) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(data.RefreshTokenHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	session := models.Session{
		UserID:           data.UserID,
		JTI:              data.JTI,
		RefreshTokenHash: string(hashed),
		UserAgent:        data.UserAgent,
		IP:               data.IP,
		CreatedAt:        time.Now(),
		RefreshExpiresIn: time.Now().Add(time.Duration(s.config.RefreshTokenExpirySEC) * time.Second),
	}
	return s.db.Create(&session).Error
}

func (s *SessionService) DeleteSessionByID(id uint) error {
	return s.db.Delete(&models.Session{}, id).Error
}

func (s *SessionService) FindSessionByJTI(jti string) (*models.Session, error) {
	var session models.Session
	if err := s.db.Where("jti = ?", jti).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *SessionService) DeleteSessionsByIPAgent(userID string, ip string, userAgent string) error {
	return s.db.Where("user_id = ? AND ip = ? AND user_agent = ?", userID, ip, userAgent).Delete(&models.Session{}).Error
}
