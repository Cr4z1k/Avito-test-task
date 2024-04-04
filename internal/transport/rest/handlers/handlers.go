package handlers

import (
	"github.com/Cr4z1k/Avito-test-task/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()

	mwAdm := r.Group("", h.identifyAdmin)
	{
		banner := mwAdm.Group("/banner")
		{
			banner.GET("", h.GetBannerWithFilter)
			banner.POST("", h.CreateBanner)
			banner.PATCH("/{id}", h.UpdateBanner)
			banner.DELETE("/{id}", h.DeleteBanner)
		}
	}

	mwUsr := r.Group("", h.identifyUser)
	{
		banner := mwUsr.Group("/user_banner")
		{
			banner.GET("", h.GetBanner)
		}
	}

	return r
}
