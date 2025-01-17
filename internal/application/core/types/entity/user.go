package entity

import (
	"time"
)

type User struct {
	ID         string     `json:"id"`         // Primary key
	FirstName  string     `json:"firstName"`  // User's first name
	LastName   string     `json:"lastName"`   // User's last name
	Email      string     `json:"email"`      // User's email (unique)
	Password   string     `json:"-"`          // Hashed password (excluded from JSON responses)
	IsVerified bool       `json:"isVerified"` // Email verification status
	Role       string     `json:"role"`       // User role (e.g., "user", "admin")
	IsLocked   bool       `json:"isLocked"`   // Whether the account is locked
	CreatedAt  time.Time  `json:"createdAt"`  // Record creation timestamp
	UpdatedAt  time.Time  `json:"updatedAt"`  // Last update timestamp
	DeletedAt  *time.Time `json:"deletedAt"`  // Soft delete timestamp (nullable)
}
