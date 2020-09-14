package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"golang.org/x/crypto/bcrypt"

	"github.com/igomonov88/users/internal/platform/auth"
)

var (
	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrAuthenticationFailure = errors.New("Authentication failed")

	// ErrNotFound is used when a specific User is requested but does not exist.
	ErrNotFound = errors.New("User not found")

	// ErrEmailAlreadyExist is used when User with such email is already exist.
	ErrEmailAlreadyExist = errors.New("Email already exist")

	// ErrUserNameAlreadyExist is used when user with such UserName is already
	// exist.
	ErrUserNameAlreadyExist = errors.New("User Name already exist")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidUserID = errors.New("ID is not in its proper form")
)

// claimsDuration represents time which our token will be valid
var claimsDuration = time.Hour

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims value representing this user. The claims can be
// used to generate a token for future authentication.
func Authenticate(ctx context.Context, db *sqlx.DB, now time.Time, email, password string) (auth.Claims, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.Authenticate")
	defer span.End()

	const q = `SELECT * FROM users WHERE email = $1;`

	var u User

	if err := db.GetContext(ctx, &u, q, email); err != nil {

		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err == sql.ErrNoRows {
			return auth.Claims{}, ErrAuthenticationFailure
		}

		return auth.Claims{}, errors.Wrap(err, "selecting single user")
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.NewClaims(u.ID, now, claimsDuration)
	return claims, nil
}

// Create user with provided info in database.
func Create(ctx context.Context, db *sqlx.DB, email, userName, avatar, password string) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.Create")
	defer span.End()

	const q = `INSERT INTO users (
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

	_, dbErr := db.NamedExec(q, u)
	if dbErr != nil {
		return nil, errors.Wrap(constraintError(dbErr), "inserting user")
	}

	return &u, nil
}

// Delete removes a user from the database.
func Delete(ctx context.Context, db *sqlx.DB, userID string) error {
	ctx, span := trace.StartSpan(ctx, "internal.user.Delete")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUserID
	}

	const q = `DELETE from users WHERE user_id = $1;`

	if _, err := db.ExecContext(ctx, q, userID); err != nil {
		return errors.Wrapf(err, "deleting user %s", userID)
	}

	return nil
}

// DeleteAvatar removes an avatar of given user from the database.
func DeleteAvatar(ctx context.Context, db *sqlx.DB, userID string) error {
	ctx, span := trace.StartSpan(ctx, "internal.user.DeleteAvatar")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUserID
	}

	const q = `UPDATE users SET avatar = '' WHERE user_id = $1;`

	if _, err := db.ExecContext(ctx, q, userID); err != nil {
		return errors.Wrapf(err, "deleting avatar userID %q", userID)
	}

	return nil
}

// DoesEmailExist returns info about existing email in database.
func DoesEmailExist(ctx context.Context, db *sqlx.DB, email string) (bool, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.DoesEmailExist")
	defer span.End()

	var exists bool
	const q = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);`

	err := db.GetContext(ctx, &exists, q, email)
	if err != nil {
		return exists, errors.Wrapf(err, "selecting email exists %q", email)
	}

	return exists, err
}

// DoesUserNameExist returns info about existing user name in database.
func DoesUserNameExist(ctx context.Context, db *sqlx.DB, userName string) (bool, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.DoesUserNameExist")
	defer span.End()

	var exist bool
	const q = `SELECT EXISTS(SELECT 1 FROM users WHERE user_name = $1);`

	err := db.GetContext(ctx, &exist, q, userName)
	if err != nil {
		return exist, errors.Wrapf(err, "selecting user name exists %q", userName)
	}

	return exist, err
}

// Retrieve gets the specified user from the database.
func Retrieve(ctx context.Context, db *sqlx.DB, userID string) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.Retrieve")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return nil, ErrInvalidUserID
	}

	const q = `SELECT * FROM users where user_id = $1;`
	var u User
	if err := db.GetContext(ctx, &u, q, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrapf(err, "selecting user %q", userID)
	}

	return &u, nil
}

// RetrieveByEmail gets the specified user from the database by email.
func RetrieveByEmail(ctx context.Context, db *sqlx.DB, email string) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.RetrieveByEmail")
	defer span.End()

	const q = `SELECT * FROM users WHERE email=$1;`
	var u User

	if err := db.GetContext(ctx, &u, q, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrapf(err, "selecting user %q", email)
	}

	return &u, nil
}

// RetrieveByUserName gets the specified user from the database by user name.
func RetrieveByUserName(ctx context.Context, db *sqlx.DB, userName string)(*User, error) {
	ctx, span := trace.StartSpan(ctx, "internal.user.RetrieveByUserName")
	defer span.End()

	const q = `SELECT * FROM users WHERE user_name=$1;`
	var u User

	if err := db.GetContext(ctx, &u, q, userName); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrapf(err, "selecting user %q", userName)
	}

	return &u, nil
}

// Update replaces a user document in the database.
func Update(ctx context.Context, db *sqlx.DB, userID, userName, email string) error {
	ctx, span := trace.StartSpan(ctx, "internal.user.Update")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUserID
	}

	const q = `UPDATE users SET user_name = $2, email = $3 WHERE 
    user_id = $1;`

	if _, err := db.ExecContext(ctx, q, userID, userName, email); err != nil {
		return constraintError(err)
	}

	return nil
}

// Update replaces a user avatar in the database.
func UpdateAvatar(ctx context.Context, db *sqlx.DB, userID, avatar string) error {
	ctx, span := trace.StartSpan(ctx, "internal.user.UpdateAvatar")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUserID
	}

	const q = `UPDATE users SET avatar = $2 WHERE user_id = $1;`

	if _, err := db.ExecContext(ctx, q, userID, avatar); err != nil {
		return errors.Wrapf(err, "updating avatar %q", userID)
	}

	return nil
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
