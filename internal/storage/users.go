package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is used when a specific User is requested but does not exist.
	ErrNotFound = errors.New("User not found")

	ErrEmailAlreadyExist = errors.New("Email already exist")

	ErrUserNameAlreadyExist = errors.New("UserName already exist")
)

func Create(ctx context.Context, db *sqlx.DB, email, userName, avatar, password string) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.Create")
	defer span.End()

	const query = `INSERT INTO users (
	user_id, user_name, email, password_hash, avatar, created_at, updated_at, 
	deleted_at) VALUES (:user_id, :user_name, :email, :password_hash, :avatar, 
	:created_at, :updated_at, :deleted_at);`

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u := User{
		ID:           uuid.New().String(),
		Name:         userName,
		Email:        email,
		PasswordHash: hash,
		Avatar:       avatar,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Time{},
		DeletedAt:    time.Time{},
	}

	_, dbErr := db.NamedExec(query, u)
	if dbErr != nil {
		return nil, errors.Wrap(constraintError(dbErr), "inserting user")
	}

	return &u, nil
}

func constraintError(err error) error {
	const UniqueViolationCode = "23505"
	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == UniqueViolationCode {
			switch pqErr.Constraint {
			case "email_idx":
				return ErrEmailAlreadyExist
			case "user_name_idx":
				return ErrUserNameAlreadyExist
			}
		}
	}
	return err
}
