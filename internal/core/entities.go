package core

type BannerContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

type BannerWithFilters struct {
	BannerID  int           `json:"id"`
	TagIDs    []int64       `json:"tag_ids"`
	FeatureID int           `json:"feature_id"`
	Content   BannerContent `json:"content"`
	IsActive  bool          `json:"is_active"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}
