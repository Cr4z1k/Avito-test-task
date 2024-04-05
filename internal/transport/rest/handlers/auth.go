package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetToken(c *gin.Context) {
	isAdmin := c.Param("isAdmin")

	var token string
	var err error

	if isAdmin == "1" {
		token, err = h.s.Auth.GetToken(true)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	} else {
		token, err = h.s.Auth.GetToken(false)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
