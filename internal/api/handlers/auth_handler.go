package handlers

import (
	"net/http"
	"test-task/internal/dto"
	"test-task/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// @Tags         auth
// @Summary      Вход по GUID
// @Description  Аутентификация пользователя по GUID и выдача токенов.
// @Param        guid path string true "GUID пользователя"
// @Success      200 {object} dto.TokenResponse
// @Failure      400 "GUID не передан или некорректный GUID"
// @Failure      500 "Внутренняя ошибка сервера"
// @Router       /api/auth/login/{guid} [get]
func (h *AuthHandler) Login(c *gin.Context) {

	guid := c.Param("guid")
	if guid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GUID required"})
		return
	}

	_, err := uuid.Parse(guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid GUID"})
		return
	}

	response, err := h.authService.Login(guid, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Tags         auth
// @Summary      Выход
// @Description  Выход пользователя и удаление сессии
// @Param        body body dto.LogoutRequest true "Access token для выхода"
// @Success      200 "Вы успешно вышли из системы"
// @Failure      400 "access_token не передан в теле запроса"
// @Failure      500 "Внутренняя ошибка сервера"
// @Router       /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.AccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access_token required"})
		return
	}

	err := h.authService.Logout(req.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// @Tags         auth
// @Summary      Обновление access токена
// @Description  Обновляет JWT access токен с помощью refresh токена
// @Param        body body dto.RefreshTokenRequest true "Запрос на обновление токена"
// @Success      200 {object} dto.TokenResponse
// @Failure      400 "Некорректный запрос"
// @Failure      401 "Неверный access/refresh токен, пользователь не найден или не совпадает User-Agent"
// @Failure      500 "Внутренняя ошибка сервера"
// @Router       /api/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	response, err := h.authService.RefreshTokens(req.AccessToken, req.RefreshToken, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
