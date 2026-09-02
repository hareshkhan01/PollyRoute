package domain

import "time"

type User struct {
	ID           string    // Can not be null
	Name         string    // Can not be null
	Email        string    // Can not be null
	PasswordHash string    // Can not be null
	CreatedAt    time.Time // Can not be null

	RefreshToken          *string    // Can be null
	RefreshTokenExpiresAt *time.Time //Can be null
}
