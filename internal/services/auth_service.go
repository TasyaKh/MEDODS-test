package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"test-task/internal/config"
	"test-task/internal/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	config         *config.Config
	db             *gorm.DB
	sessionService *SessionService
	userService    *UserService
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthService(cfg *config.Config, db *gorm.DB, sessionService *SessionService, userService *UserService) *AuthService {
	return &AuthService{
		config:         cfg,
		db:             db,
		sessionService: sessionService,
		userService:    userService,
	}
}

func (s *AuthService) GenerateRefreshToken() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func (s *AuthService) GenerateJWTToken(userID string) (*dto.JWTTokenData, error) {
	expirationTime := time.Now().Add(time.Duration(s.config.JWTAccessExpirySEC) * time.Second)
	jti := uuid.NewString()
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &dto.JWTTokenData{
		Token:     tokenString,
		JTI:       jti,
		ExpiresIn: s.config.JWTAccessExpirySEC,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *AuthService) ValidateTokenIgnoreExpiry(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

func (s *AuthService) SendWebhookNotification(userId string, ipAddress string) {

	if s.config.WebhookURL == "" {
		return
	}

	payload := dto.WebhookPayload{
		IPAddress: ipAddress,
		Event:     "ip_change_detected",
	}

	go func() {
		body, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("failed to marshal webhook payload: %v\n", err)
			return
		}

		resp, err := http.Post(s.config.WebhookURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("failed to send webhook notification: %v\n", err)
			return
		}
		defer resp.Body.Close()
	}()
}

func (s *AuthService) Login(userID, userAgent, ip string) (*dto.TokenResponse, error) {

	user, err := s.userService.GetUserByID(userID)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	refreshToken := s.GenerateRefreshToken()

	jwtTokenResponse, err := s.GenerateJWTToken(userID)

	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	err = s.sessionService.DeleteSessionsByIPAgent(user.ID.String(), ip, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to delete sessions: %w", err)
	}

	err = s.sessionService.AddSession(dto.SessionCreate{
		UserID:           user.ID,
		RefreshTokenHash: refreshToken,
		IP:               ip,
		UserAgent:        userAgent,
		JTI:              jwtTokenResponse.JTI,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session")
	}

	return &dto.TokenResponse{
		AccessToken:  jwtTokenResponse.Token,
		RefreshToken: refreshToken,
		ExpiresIn:    jwtTokenResponse.ExpiresIn,
	}, nil
}

// compare refresh token with hash
func (s *AuthService) CompareRefreshToken(hashedToken, token string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
}

func (s *AuthService) Logout(token string) error {
	claims, err := s.ValidateTokenIgnoreExpiry(token)
	if err != nil {
		return err
	}

	session, err := s.sessionService.FindSessionByJTI(claims.ID)
	if err != nil || session == nil {
		return fmt.Errorf("session not found")
	}

	if err := s.sessionService.DeleteSessionByID(session.ID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *AuthService) RefreshTokens(accessToken, refreshToken, userAgent, ip string) (*dto.TokenResponse, error) {
	claims, err := s.ValidateTokenIgnoreExpiry(accessToken)
	if err != nil {
		return nil, err
	}

	session, err := s.sessionService.FindSessionByJTI(claims.ID)
	if err != nil || session == nil {
		return nil, fmt.Errorf("session not found")
	}

	// refresh expired
	if time.Now().After(session.RefreshExpiresIn) {
		_ = s.sessionService.DeleteSessionByID(session.ID)
		return nil, fmt.Errorf("refresh token expired")
	}

	if err := s.CompareRefreshToken(session.RefreshTokenHash, refreshToken); err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// check User-Agent
	if session.UserAgent != userAgent {
		_ = s.sessionService.DeleteSessionByID(session.ID)
		return nil, fmt.Errorf("user agent mismatch: deauthorized")
	}

	// check IP and send webhook if changed
	if session.IP != ip {
		s.SendWebhookNotification(claims.UserID, ip)
	}

	err = s.sessionService.DeleteSessionByID(session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete old session: %w", err)
	}

	loginResponse, err := s.Login(claims.UserID, userAgent, ip)

	if err != nil {
		return nil, fmt.Errorf("failed to refresh tokens: %w", err)
	}
	return loginResponse, nil

}
