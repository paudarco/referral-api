package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paudarco/referral-api/models"
)

// CreateReferralCode создает новый реферальный код
func (h *Handler) CreateReferralCode(c *gin.Context) {
	var req models.CreateReferralCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "wrong date format")
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	code, err := h.services.CreateReferralCode(c.Request.Context(), userID, req.ExpiresAt)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, code)
}

// GetReferrals получает список рефералов пользователя
func (h *Handler) GetReferrals(c *gin.Context) {
	userID := c.GetInt("user_id")
	referrals, err := h.services.GetReferrals(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	c.JSON(http.StatusOK, referrals)
}

func (h *Handler) GetReferralCodeByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		newErrorResponse(c, http.StatusBadRequest, "missing email parameter")
		return
	}
	code, err := h.services.GetReferralCodeByEmail(c.Request.Context(), email)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, code)
}

func (h *Handler) DeleteReferralCode(c *gin.Context) {
	userID := c.GetInt("user_id")

	err := h.services.DeleteReferralCode(c.Request.Context(), userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
