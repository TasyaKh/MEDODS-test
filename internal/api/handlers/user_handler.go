package handlers

import (
	"net/http"
	"test-task/internal/models"
	"test-task/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Tags         user
// @Security     BearerAuth
// @Summary      Получить текущего пользователя
// @Description  Возвращает текущего аутентифицированного пользователя.
// @Success 200 "Успешный ответ: 'id' пользователя"
// @Failure      401 "Пользователь не авторизован"
// @Failure      500 "Внутренняя ошибка сервера"
// @Router       /api/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	rawUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := rawUser.(*models.User)
	c.JSON(http.StatusOK, gin.H{
		"id": user.ID,
	})
}
