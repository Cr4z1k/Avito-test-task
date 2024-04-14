package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetBanner(c *gin.Context) {
	isAdmin, ok := c.Get("isAdmin")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unable to get value by key from context"})
		return
	}

	isAdminBool, ok := isAdmin.(bool)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to perform type conversion"})
		return
	}

	tagIDStr := c.Query("tag_id")

	tagID, err := strconv.ParseUint(tagIDStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	featureIDStr := c.Query("feature_id")

	featureID, err := strconv.ParseUint(featureIDStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var useLastRevision bool

	useLastRevisionStr := c.Query("use_last_revision")

	if useLastRevisionStr == "" {
		useLastRevision = false
	} else {
		useLastRevision, err = strconv.ParseBool(useLastRevisionStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	bannerContent, err := h.s.GetBanner(tagID, featureID, useLastRevision, isAdminBool)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if bannerContent == (core.BannerContent{}) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, bannerContent)
}

func (h *Handler) GetBannerWithFilter(c *gin.Context) {
	var (
		featureID *int
		tagID     *int
		limit     int = 100
		offset    int = 0
		err       error
	)

	featureIDString := c.Query("feature_id")

	if featureIDString != "" {
		featureIDconv, err := strconv.Atoi(featureIDString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if featureIDconv < 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID cannot be less than 0"})
			return
		}

		featureID = &featureIDconv
	}

	tagIDString := c.Query("tag_id")

	if tagIDString != "" {
		tagIDconv, err := strconv.Atoi(tagIDString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if tagIDconv < 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID cannot be less than 0"})
			return
		}

		tagID = &tagIDconv
	}

	limitString := c.Query("limit")

	if limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	offsetString := c.Query("offset")

	if offsetString != "" {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	searchResult, err := h.s.Banner.GetBannersWithFilter(tagID, featureID, limit, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, searchResult)
}

func (h *Handler) CreateBanner(c *gin.Context) {
	var jsonInfo BannerJSON

	if err := c.ShouldBindJSON(&jsonInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(jsonInfo.TagIDs) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "slice has 0 elements"})
		return
	}

	var flag bool = false

	for _, tagID := range jsonInfo.TagIDs {
		if tagID < 0 {
			flag = true
			break
		}
	}

	if flag || jsonInfo.FeatureID < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID cannot be less than 0"})
		return
	}

	id, err := h.s.Banner.CreateBanner(jsonInfo.TagIDs, uint64(jsonInfo.FeatureID), jsonInfo.Content, jsonInfo.IsActive)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"banner_id": id})
}

func (h *Handler) UpdateBanner(c *gin.Context) {
	bannerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var jsonInfo BannerJSON

	if err = c.ShouldBindJSON(&jsonInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(jsonInfo.TagIDs) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "slice has 0 elements"})
		return
	}

	var flag bool

	for _, tagID := range jsonInfo.TagIDs {
		if tagID < 0 {
			flag = true
			break
		}
	}

	if flag || jsonInfo.FeatureID < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID cannot be less than 0"})
		return
	}

	err = h.s.UpdateBanner(bannerID, uint64(jsonInfo.FeatureID), jsonInfo.TagIDs, jsonInfo.Content, jsonInfo.IsActive)
	if err != nil {
		if strings.Contains(err.Error(), "No banner with such ID") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) DeleteBanner(c *gin.Context) {
	bannerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.s.Banner.DeleteBanner(bannerID)
	if err == sql.ErrNoRows {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
