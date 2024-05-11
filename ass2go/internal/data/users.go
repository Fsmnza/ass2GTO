package data

import (
	"ass2/internal/validator"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)
var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Fname     string    `json:"fname"`
	Sname     string    `json:"sname"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Password  password  `json:"password"`
	Activated bool      `json:"activated"`
	Version   int       `json:"version"`
}
type UserModel struct {
	DB *sql.DB
}

type password struct {
	plaintext *string
	hash      []byte
}

func (m UserModel) Insert(user *User) error {
	query := `
INSERT INTO users (fname, sname,role, email, password_hash, activated)
VALUES ($1, $2, 'user', $3, $4, false)
RETURNING id, created_at, version`
	args := []any{user.Fname, user.Sname, user.Email, user.Password.hash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
SELECT id, created_at, updated_at, fname, sname, email, role, password_hash, activated, version
FROM users
WHERE email = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Fname,
		&user.Sname,
		&user.Email,
		&user.Role,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
UPDATE users
SET fname = $1,sname = $2, updated_at = $3,email = $4, password_hash = $5, activated = $6, version = version + 1
WHERE id = $7 AND version = $8
RETURNING version`
	args := []any{
		user.Fname,
		user.Sname,
		user.UpdatedAt,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `
SELECT users.id, users.created_at, users.updated_at,users.fname, users.sname, users.role, users.email, users.password_hash, users.activated, users.version
FROM users
INNER JOIN tokens
ON users.id = tokens.user_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > $3`
	args := []any{tokenHash[:], tokenScope, time.Now()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Fname,
		&user.Sname,
		&user.Role,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
func (m UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
			SELECT *
			FROM users
			WHERE id = $1`

	var user User

	err := m.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Fname,
		&user.Sname,
		&user.Email,
		&user.Password.hash,
		&user.Role,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) GetAll() ([]*User, error) {
	query := `SELECT id, created_at, updated_at, fname, sname, email, password_hash, role, activated, version FROM users`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userInfos := []*User{}
	for rows.Next() {
		userInfo := &User{}
		err = rows.Scan(
			&userInfo.ID,
			&userInfo.CreatedAt,
			&userInfo.UpdatedAt,
			&userInfo.Fname,
			&userInfo.Sname,
			&userInfo.Email,
			&userInfo.Password.hash,
			&userInfo.Role,
			&userInfo.Activated,
			&userInfo.Version,
		)
		if err != nil {
			return nil, err
		}
		userInfos = append(userInfos, userInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userInfos, nil
}

func (m UserModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM users
		WHERE id = $1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Fname != "", "name", "must be provided")
	v.Check(len(user.Fname) <= 500, "name", "must not be more than 500 bytes long")
	ValidateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
func (m UserModel) FindNotActivatedAndExpired() ([]*User, error) {
	query := `
      SELECT u.*
      FROM users u
      INNER JOIN public.tokens uit on u.id = uit.user_id
      WHERE u.activated = false AND uit.expiry < now() - interval '10 s'
`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Fname, &user.Sname, &user.Email, &user.Password.hash, &user.Role, &user.Activated, &user.Version)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
