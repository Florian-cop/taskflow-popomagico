package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"taskflow-api/internal/user/domain"
)

type JWTTokenGenerator struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTTokenGenerator(secret string, ttl time.Duration) *JWTTokenGenerator {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &JWTTokenGenerator{secret: []byte(secret), ttl: ttl}
}

type claims struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	jwt.RegisteredClaims
}

func (g *JWTTokenGenerator) Generate(_ context.Context, user *domain.User) (string, error) {
	now := time.Now()
	c := claims{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(g.ttl)),
			Issuer:    "taskflow",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(g.secret)
}

func (g *JWTTokenGenerator) Validate(_ context.Context, tokenStr string) (*domain.User, error) {
	parsed, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return g.secret, nil
	})
	if err != nil {
		return nil, err
	}

	c, ok := parsed.Claims.(*claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	return &domain.User{
		ID:        c.Subject,
		Email:     c.Email,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}, nil
}
