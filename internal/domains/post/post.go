package postdomain

import "time"

type Post struct {
	ID         string    `json:"id"`
	AuthorID   string    `json:"author_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CategoryID *string   `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
