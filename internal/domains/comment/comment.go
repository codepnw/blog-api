package commentdomain

import "time"

type Comment struct {
	ID        int64
	PostID    string
	UserID    string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
