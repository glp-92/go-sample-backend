package user

import "github.com/google/uuid"

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDetailsResponse struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}
