package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateFeatures(c *gin.Context) {
	jsonInfo := struct {
		Features []int `json:"features"`
	}{}

	if err := c.ShouldBindJSON(&jsonInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(jsonInfo.Features) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "slice has 0 elements"})
		return
	}

	if err := h.s.Feature.CreateFeatures(jsonInfo.Features); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
