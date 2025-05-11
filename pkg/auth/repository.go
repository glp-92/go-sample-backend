package auth

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type AuthRepository interface {
	SaveUser(user User) error
	GetUserDetails(username string) (*User, error)
	SaveRefreshToken(refreshToken RefreshToken) error
	GetRefreshTokenFromSubject(username string) (RefreshToken, error)
	RevokeRefreshToken(userID uuid.UUID) error
}

type MySQLAuthRepository struct {
	db *sql.DB
}

func NewMySQLAuthRepository(db *sql.DB) AuthRepository {
	return &MySQLAuthRepository{db: db}
}

func (r *MySQLAuthRepository) SaveUser(user User) error {
	if user.Username == "" || user.Password == "" {
		return errors.New("user fields cannot be empty")
	}
	query := `
		INSERT INTO users (id, username, password)
		VALUES (?, ?, ?)
	`
	_, err := r.db.Exec(query, user.Id, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

func (r *MySQLAuthRepository) GetUserDetails(username string) (*User, error) {
	query := `
		SELECT id, username, password 
		FROM users
		WHERE username = ?
	`
	row := r.db.QueryRow(query, username)
	var user User
	err := row.Scan(&user.Id, &user.Username, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	return &user, nil
}

func (r *MySQLAuthRepository) SaveRefreshToken(refreshToken RefreshToken) error {
	query := `
		INSERT INTO tokens (id, user_id, refresh_token, revoked)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			refresh_token = VALUES(refresh_token),
			revoked = VALUES(revoked)
	`
	_, err := r.db.Exec(query, refreshToken.Id, refreshToken.UserId, refreshToken.RefreshToken, refreshToken.Revoked)
	if err != nil {
		return fmt.Errorf("error saving refresh token: %w", err)
	}
	return nil
}

func (r *MySQLAuthRepository) GetRefreshTokenFromSubject(username string) (RefreshToken, error) {
	query := `
		SELECT tokens.id, tokens.user_id, tokens.refresh_token, tokens.revoked
		FROM tokens
		JOIN users ON users.id = tokens.user_id
		WHERE users.username = ?
	`
	row := r.db.QueryRow(query, username)
	var refreshToken RefreshToken
	err := row.Scan(&refreshToken.Id, &refreshToken.UserId, &refreshToken.RefreshToken, &refreshToken.Revoked)
	return refreshToken, err
}

func (r *MySQLAuthRepository) RevokeRefreshToken(userId uuid.UUID) error {
	query := `UPDATE tokens SET revoked = TRUE WHERE user_id = ?`
	_, err := r.db.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error revoking refresh token for user %s: %w", userId, err)
	}
	return nil
}
