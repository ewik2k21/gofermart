package services

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"gofermart/config"
	"gofermart/internals/interfaces"
	"time"
)

type ITokenService interface {
	GenerateJwtToken(id string) (*string, time.Time, error)
}

type TokenService struct {
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *TokenService) GenerateJwtToken(id string) (*string, time.Time, error) {
	expirationTime := time.Now().UTC().Add(time.Hour * 24)

	claims := &interfaces.Claims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JwtKey))

	if err != nil {
		logrus.Errorf("jwt token not signed: %s", err)
		return nil, time.Now(), err
	}

	return &tokenString, expirationTime, nil
}
