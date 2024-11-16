package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/referral-api/internal/errors"
	"github.com/paudarco/referral-api/models"
)

func (h *Handler) SignUp(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid data format")
		return
	}

	response, err := h.services.Register(c.Request.Context(), &req)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, response.Token)
}

func (h *Handler) SignIn(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid data format")
		return
	}

	response, err := h.services.Login(c.Request.Context(), &req)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, errors.ErrInvalidCredentials.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": response.Token,
	})
}
