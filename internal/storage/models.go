package storage

import (
	"time"
)

// User represents someone with access to our system.
type User struct {
	ID           string    `db:"user_id"`
	Name         string    `db:"user_name"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	Avatar       string    `db:"avatar"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	DeletedAt    time.Time `db:"deleted_at"`
}
