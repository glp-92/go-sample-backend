package auth

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Username string
	Password string
}

type RefreshToken struct {
	Id           uuid.UUID
	RefreshToken string
	Revoked      bool
	UserId       uuid.UUID
}
