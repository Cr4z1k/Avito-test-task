package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CheckTokenIsAdmin(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if header == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenParts := strings.Split(header, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if len(tokenParts[0]) == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	isAdmin, err := h.s.Auth.ParseToken(tokenParts[1], "isAdmin")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("isAdmin", isAdmin)
}

func (h *Handler) IdentifyAdmin(c *gin.Context) {
	isAdmin, ok := c.Get("isAdmin")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	isAdminBool, ok := isAdmin.(bool)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !isAdminBool {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}
