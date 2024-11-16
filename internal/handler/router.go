package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/paudarco/referral-api/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.Default()

	// Публичные маршруты
	public := router.Group("/api/v1")
	{
		public.POST("/register", h.SignUp)
		public.POST("/login", h.SignIn)
		public.GET("/referral-code", h.GetReferralCodeByEmail)
	}

	// Защищенные маршруты
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware(h))
	{
		protected.POST("/referral-code", h.CreateReferralCode)
		protected.DELETE("/referral-code", h.DeleteReferralCode)
		protected.GET("/referrals", h.GetReferrals)
	}

	return router
}
