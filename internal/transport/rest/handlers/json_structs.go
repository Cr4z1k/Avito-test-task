package handlers

import "github.com/Cr4z1k/Avito-test-task/internal/core"

type GetBanner struct {
	TagID           int   `json:"tag_id"`
	FeatureID       int   `json:"feature_id"`
	UseLastRevision *bool `json:"use_last_revision"`
}

type CreateBanner struct {
	TagIDs    []int              `json:"tag_ids"`
	FeatureID int                `json:"feature_id"`
	Content   core.BannerContent `json:"content"`
	IsActive  bool               `json:"is_active"`
}
