package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MapClaims = jwt.MapClaims

type Service struct {
	SecretKey   string
	ExpiredTime time.Duration
}

func New() Service {
	return Service{
		SecretKey:   os.Getenv("JWT_SECRET_KEY"),
		ExpiredTime: 1,
	}
}

type jwtCustomClaim struct {
	jwt.StandardClaims
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (s Service) GenerateToken(
	id, email, role string,
) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&jwtCustomClaim{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * s.ExpiredTime).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			ID:    id,
			Email: email,
			Role:  role,
		},
	)

	t, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s Service) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.SecretKey), nil
	})
}
