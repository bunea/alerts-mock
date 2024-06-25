package models

import "time"

type Entry struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" faker:"sentence"`
	PublishedAt time.Time `json:"published_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Href        string    `json:"href" faker:"url"`
	ImageURL    string    `json:"image_url"`
	Content     string    `json:"content" faker:"paragraph"`
	AuthorName  string    `json:"author_name" faker:"name"`
}

type Feed struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	UpdatedAt   time.Time `json:"updated_at"`
	Entries     []*Entry  `json:"entries"`
	UpdateEvery Duration  `json:"update_every"`
}
