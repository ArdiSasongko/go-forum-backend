package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/ArdiSasongko/go-forum-backend/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type ClaimsToken struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsValid  bool   `json:"is_valid"`
	jwt.RegisteredClaims
}

var MapToken = map[string]time.Duration{
	"token":         30 * time.Minute,
	"refresh_token": 60 * 24 * 10 * time.Minute,
}

var secretKey = []byte(env.GetEnv("JWT_SECRET", ""))

func GenerateToken(ctx context.Context, claims ClaimsToken, tokenType string) (string, error) {
	claimsToken := ClaimsToken{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
		IsValid:  claims.IsValid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(MapToken[tokenType]).UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsToken)
	result, err := token.SignedString(secretKey)
	if err != nil {
		logrus.WithField("jwt token", err.Error()).Error(err.Error())
		return result, fmt.Errorf("failed generate jwt token: %v", err)
	}
	return result, nil
}

func ValidateToken(ctx context.Context, token string) (*ClaimsToken, error) {
	var (
		claimToken *ClaimsToken
		ok         bool
	)

	jwtToken, err := jwt.ParseWithClaims(token, &ClaimsToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.WithField("jwt token", t.Header["alg"]).Error(t.Header["alg"])
			return nil, fmt.Errorf("failed validate method jwt : %v", t.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		logrus.WithField("jwt token", err.Error()).Error(err.Error())
		return nil, fmt.Errorf("failed parse token : %v", err)
	}

	if claimToken, ok = jwtToken.Claims.(*ClaimsToken); !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("token invalid : %v", err)
	}

	return claimToken, nil
}

func ValidateRefreshToken(ctx context.Context, token string) (*ClaimsToken, error) {
	var (
		claimToken *ClaimsToken
		ok         bool
	)

	jwtToken, err := jwt.ParseWithClaims(token, &ClaimsToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.WithField("jwt token", t.Header["alg"]).Error(t.Header["alg"])
			return nil, fmt.Errorf("failed validate method jwt : %v", t.Header["alg"])
		}
		return secretKey, nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		logrus.WithField("jwt token", err.Error()).Error(err.Error())
		return nil, fmt.Errorf("failed parse token : %v", err)
	}

	if claimToken, ok = jwtToken.Claims.(*ClaimsToken); !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("token invalid : %v", err)
	}

	return claimToken, nil
}
