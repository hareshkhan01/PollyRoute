package request

// Use for registration request
type NewUser struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	RawPassword string `json:"password" binding:"required"`
}

// Use for login request
type OldUser struct {
	Email       string `json:"email" binding:"required"`
	RawPassword string `json:"password" binding:"required"`
}
