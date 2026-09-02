package domain

import "time"

type User struct {
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt time.Time
}
