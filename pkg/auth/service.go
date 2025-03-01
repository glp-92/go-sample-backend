package auth

import (
	"errors"
	"fmt"
	"fullstackcms/backend/pkg/auth/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccessTokenClaims struct {
	Expires *jwt.NumericDate `json:"expires"`
	Issued  *jwt.NumericDate `json:"issued"`
	Subject string           `json:"sub"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	AccessTokenClaims
	Agent string `json:"agent"`
}

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

func (s *AuthService) ValidateUser(request dto.LoginRequest) (*User, error) {
	user, err := s.repo.GetUserDetails(request.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("User not found Error")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) CreateTokens(request dto.LoginRequest, userAgent string, user *User) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expires": jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		"issued":  jwt.NewNumericDate(time.Now()),
		"sub":     user.Username,
	})
	signedAccessToken, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshTokenClaims{
		AccessTokenClaims: AccessTokenClaims{
			Expires: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			Issued:  jwt.NewNumericDate(time.Now()),
			Subject: user.Username,
		},
		Agent: userAgent,
	})
	signedRefreshToken, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}
	newRefreshToken := RefreshToken{
		Id:           uuid.New(),
		RefreshToken: signedRefreshToken,
		UserId:       user.Id,
	}
	err = s.repo.SaveRefreshToken(newRefreshToken)
	if err != nil {
		return "", "", err
	}
	return signedAccessToken, newRefreshToken.RefreshToken, err
}

func (s *AuthService) RefreshToken(refreshToken string) error {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	claims := token.Claims.(*RefreshTokenClaims)
	fmt.Println(claims.Subject)
	return err
}
