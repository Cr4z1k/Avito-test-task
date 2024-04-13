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

	r.GET("/get_token/:isAdmin", h.GetToken)

	mwToken := r.Group("", h.CheckTokenIsAdmin)
	{
		mwAdm := mwToken.Group("", h.IdentifyAdmin)
		{
			banner := mwAdm.Group("/banner")
			{
				banner.GET("", h.GetBannerWithFilter)
				banner.POST("", h.CreateBanner)
				banner.PATCH("/:id", h.UpdateBanner)
				banner.DELETE("/:id", h.DeleteBanner)
			}

			tag := mwAdm.Group("/tag")
			{
				tag.POST("", h.CreateTags)
			}

			feature := mwAdm.Group("/feature")
			{
				feature.POST("", h.CreateFeatures)
			}
		}

		mwToken.GET("/user_banner", h.GetBanner)
	}

	return r
}
