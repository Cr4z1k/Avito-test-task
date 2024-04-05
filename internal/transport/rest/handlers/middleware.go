package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) checkToken(c *gin.Context) {
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

	c.Set("tokenPart", tokenParts[1])
}

func (h *Handler) identifyAdmin(c *gin.Context) {
	tokenPart, ok := c.Get("tokenPart")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	isAdmin, err := h.s.Auth.ParseToken(tokenPart.(string), "isAdmin")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	isAdminBool, ok := isAdmin.(bool)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": ok})
		return
	}

	if !isAdminBool {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}
