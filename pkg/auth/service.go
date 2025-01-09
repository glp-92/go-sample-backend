package auth

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(request RegisterRequest) error {
	encodedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := User{
		Id:       uuid.New(),
		Username: request.Username,
		Password: string(encodedPassword),
	}
	return s.repo.Save(newUser)
}

func (s *UserService) ValidateUser(request LoginRequest) error {
	user, err := s.repo.GetUserDetails(request.Username)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("User not found Error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return err
	}
	return nil
}
