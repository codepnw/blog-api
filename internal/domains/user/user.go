package userdomain

import "time"

type User struct {
	ID           string
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
