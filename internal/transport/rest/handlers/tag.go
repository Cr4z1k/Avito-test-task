package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateTags(c *gin.Context) {
	jsonInfo := struct {
		Tags []int `json:"tags"`
	}{}

	if err := c.ShouldBindJSON(&jsonInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(jsonInfo.Tags) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "slice has 0 elements"})
		return
	}

	if err := h.s.Tag.CreateTags(jsonInfo.Tags); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
