package utils

import (
	"time"

	model "github.com/MdZunaed/go_hrms/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("jwt-secret-hrms")

func GenerateToken(user *model.User) (string, error) {
	method := jwt.SigningMethodHS256
	claims := jwt.MapClaims{
		"userId":   user.Id,
		"username": user.Email,
		"exp":      time.Now().Add(time.Hour * 180).Unix(),
	}
	token, err := jwt.NewWithClaims(method, claims).SignedString(JWTSecret)

	if err != nil {
		return "", err
	}

	return token, nil
}
