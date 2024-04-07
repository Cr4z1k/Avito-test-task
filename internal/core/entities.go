package core

type BannerContent struct {
	Title string
	Text  string
	Url   string
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
