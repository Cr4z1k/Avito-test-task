package core

type BannerContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

type BannerWithFilters struct {
	BannerID  uint64
	TagIDs    []uint64
	FeatureID uint64
	Content   BannerContent
	IsActive  bool
	CreatedAt string
	UpdatedAt string
}
