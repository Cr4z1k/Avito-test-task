package handlers

import "github.com/Cr4z1k/Avito-test-task/internal/core"

type BannerJSON struct {
	TagIDs    []int              `json:"tag_ids"`
	FeatureID int                `json:"feature_id"`
	Content   core.BannerContent `json:"content"`
	IsActive  bool               `json:"is_active"`
}
