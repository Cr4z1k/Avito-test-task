package handlers

type GetBanner struct {
	TagID           int   `json:"tag_id"`
	FeatureID       int   `json:"feature_id"`
	UseLastRevision *bool `json:"use_last_revision"`
}
