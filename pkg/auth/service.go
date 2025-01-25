package auth

import (
	"errors"
	"fullstackcms/backend/pkg/auth/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

var secretKey = []byte("secret-key")

func (s *AuthService) CreateUser(request dto.RegisterRequest) error {
	encodedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := User{
		Id:       uuid.New(),
		Username: request.Username,
		Password: string(encodedPassword),
	}
	return s.repo.SaveUser(newUser)
}

func (s *AuthService) ValidateUser(request dto.LoginRequest) error {
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

func (s *AuthService) CreateToken(request dto.LoginRequest, userAgent string) (dto.LoginResponse, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expires": jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"issued":  jwt.NewNumericDate(time.Now()),
		"sub":     request.Username,
	})
	signedAccessToken, err := accessToken.SignedString(secretKey)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expires": jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		"issued":  jwt.NewNumericDate(time.Now()),
		"sub":     request.Username,
		"agent":   userAgent,
	})
	signedRefreshToken, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	response := dto.LoginResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}
	return response, err
}

func (s *AuthService) ValidateToken() {}
