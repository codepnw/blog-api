package postdomain

import "time"

type Post struct {
	ID         string
	AuthorID   string
	Title      string
	Content    string
	CategoryID *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
