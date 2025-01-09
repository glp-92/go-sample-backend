package user

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type UserRepository interface {
	Save(user User) error
	FindByID(id uuid.UUID) (*User, error)
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

func (r *MySQLUserRepository) FindByID(id uuid.UUID) (*User, error) {
	query := `
        SELECT id, username
        FROM users
        WHERE id = ?`
	row := r.db.QueryRow(query, id)
	var user User
	err := row.Scan(&user.Id, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
