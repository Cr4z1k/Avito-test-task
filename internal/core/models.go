package core

type Banner struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}
