package auth

import (
	"database/sql"
	"errors"
)

type UserRepository interface {
	Save(user User) error
	GetUserDetails(username string) (*User, error)
}

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) UserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Save(user User) error {
	query := `
		INSERT INTO users (id, username, password)
		VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, user.Id, user.Username, user.Password)
	return err
}

func (r *MySQLUserRepository) GetUserDetails(username string) (*User, error) {
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
