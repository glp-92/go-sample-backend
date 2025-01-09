package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(request CreateUserRequest) error {
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

func (s *UserService) FindUserById(id uuid.UUID) (UserDetailsResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return UserDetailsResponse{}, err
	}
	if user == nil {
		return UserDetailsResponse{}, nil
	}
	response := UserDetailsResponse{
		Id:       user.Id,
		Username: user.Username,
	}
	return response, nil
}
