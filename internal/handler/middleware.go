package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/referral-api/internal/errors"
)

func authMiddleware(h *Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			newErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		// Bearer token
		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			newErrorResponse(c, http.StatusUnauthorized, "Invalid token format")
			return
		}

		claimsId, err := h.services.ValidateToken(parts[1])
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		// Добавляем ID пользователя в контекст
		c.Set("user_id", claimsId)
		c.Next()
	}
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("user_id")
	if !ok {
		return 0, errors.ErrUserIdNotFound
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.ErrInvalidTypeId
	}

	return idInt, nil
}
