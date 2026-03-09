package main

import "time"

// UserV1 represents the v1 API user response.
// This is the original format that existing clients depend on.
type UserV1 struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	AvatarURL string    `json:"avatar_url"`
	Bio       string    `json:"bio"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserV2 represents the v2 API user response.
// BUG: Several fields from v1 are missing (avatar_url, bio, location).
// The v2 struct was created by copying v1 and adding new fields,
// but some fields were accidentally deleted during the refactor.
// Clients that upgraded from v1 to v2 will silently lose these fields.
type UserV2 struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	// avatar_url, bio, and location were here in v1 but dropped in v2
	Department  string    `json:"department"`
	PhoneNumber string    `json:"phone_number"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InternalUser is the full internal representation.
type InternalUser struct {
	ID          string
	Name        string
	Email       string
	Role        string
	AvatarURL   string
	Bio         string
	Location    string
	Department  string
	PhoneNumber string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ToV1 converts internal user to v1 API response.
func (u *InternalUser) ToV1() UserV1 {
	return UserV1{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		AvatarURL: u.AvatarURL,
		Bio:       u.Bio,
		Location:  u.Location,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToV2 converts internal user to v2 API response.
// BUG: Does not include avatar_url, bio, or location because
// UserV2 struct doesn't have those fields.
func (u *InternalUser) ToV2() UserV2 {
	return UserV2{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Role:        u.Role,
		Department:  u.Department,
		PhoneNumber: u.PhoneNumber,
		IsActive:    u.IsActive,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
