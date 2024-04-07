package core

import "time"

type Banner struct {
	ID        uint64    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Text      string    `json:"text" db:"text"`
	Url       string    `json:"url" db:"url"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Feature struct {
	ID uint64 `json:"id" db:"id"`
}

type Tag struct {
	ID uint64 `json:"id" db:"id"`
}
