package auth

import (
	"errors"
	"fullstackcms/backend/configs"
	"fullstackcms/backend/pkg/auth/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccessTokenClaims struct {
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	Agent string `json:"agent"`
}

type AuthService struct {
	repo    AuthRepository
	authCfg configs.AuthConfig
}

func NewAuthService(repo AuthRepository, authCfg configs.AuthConfig) *AuthService {
	return &AuthService{repo: repo, authCfg: authCfg}
}

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

func (s *AuthService) CreateTokens(userAgent string, user *User) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.authCfg.JWTAccessTokenExpiration) * time.Minute)),
			Subject:   user.Username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	signedAccessToken, err := accessToken.SignedString(s.authCfg.JWTSignKey)
	if err != nil {
		return "", "", err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.authCfg.JWTRefreshTokenExpiration) * time.Minute)),
			Subject:   user.Username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Agent: userAgent,
	})
	signedRefreshToken, err := refreshToken.SignedString(s.authCfg.JWTSignKey)
	if err != nil {
		return "", "", err
	}
	newRefreshToken := RefreshToken{
		Id:           uuid.New(),
		RefreshToken: signedRefreshToken,
		UserId:       user.Id,
		Revoked:      false,
	}
	err = s.repo.SaveRefreshToken(newRefreshToken)
	if err != nil {
		return "", "", err
	}
	return signedAccessToken, newRefreshToken.RefreshToken, err
}

func (s *AuthService) RefreshToken(userAgent string, refreshToken string) (string, string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.authCfg.JWTSignKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return "", "", err
	}
	claims := token.Claims.(*RefreshTokenClaims)
	if claims.Agent != userAgent {
		return "", "", errors.New("invalid user agent")
	}
	storedRefreshToken, err := s.repo.GetRefreshTokenFromSubject(claims.Subject)
	if err != nil {
		return "", "", err
	}
	if (storedRefreshToken.Revoked) || (storedRefreshToken.RefreshToken != refreshToken) {
		return "", "", errors.New("invalid refresh token")
	}
	user, err := s.repo.GetUserDetails(claims.Subject)
	if err != nil {
		return "", "", err
	}
	var accessToken string
	accessToken, refreshToken, err = s.CreateTokens(userAgent, user)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, err
}

func (s *AuthService) ValidateTokenFromUser(accessToken string) (*User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.authCfg.JWTSignKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	user, err := s.repo.GetUserDetails(claims.Subject)
	return user, err
}

func (s *AuthService) ValidateExpiredTokenFromUser(accessToken string) (*User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.authCfg.JWTSignKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	user, err := s.repo.GetUserDetails(claims.Subject)
	return user, err
}

func (s *AuthService) Logout(userId uuid.UUID) error {
	err := s.repo.RevokeRefreshToken(userId)
	return err
}
