package auth

import (
	"database/sql"
	"errors"
	"fmt"
)

type AuthRepository interface {
	SaveUser(user User) error
	GetUserDetails(username string) (*User, error)
	SaveRefreshToken(refreshToken RefreshToken) error
}

type MySQLAuthRepository struct {
	db *sql.DB
}

func NewMySQLAuthRepository(db *sql.DB) AuthRepository {
	return &MySQLAuthRepository{db: db}
}

func (r *MySQLAuthRepository) SaveUser(user User) error {
	query := `
		INSERT INTO users (id, username, password)
		VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, user.Id, user.Username, user.Password)
	return err
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
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MySQLAuthRepository) SaveRefreshToken(refreshToken RefreshToken) error {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM tokens WHERE user_id = ?", refreshToken.UserId).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking existing token: %v", err)
	}
	if count > 0 {
		_, err = r.db.Exec(`
			UPDATE tokens 
			SET refresh_token = ?
			WHERE user_id = ?`,
			refreshToken.RefreshToken, refreshToken.UserId,
		)
		if err != nil {
			return fmt.Errorf("error updating refresh token: %v", err)
		}
	} else {
		_, err = r.db.Exec(`
			INSERT INTO tokens (id, refresh_token, user_id) 
			VALUES (?, ?, ?, ?)`,
			refreshToken.Id, refreshToken.RefreshToken, refreshToken.UserId,
		)
		if err != nil {
			return fmt.Errorf("error inserting refresh token: %v", err)
		}
	}
	return nil
}
