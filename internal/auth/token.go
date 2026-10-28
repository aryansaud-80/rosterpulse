package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type Claims struct {
	UserID       string `json:"user_id"`
	SessionID    string `json:"sid"`
	TokenVersion int    `json:"token_version"`
	TokenType    string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secretKey            string
	issuer               string
	audience             string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewTokenManager(secretKey, issuer, audience string, accessTokenDuration, refreshTokenDuration time.Duration) *TokenManager {
	return &TokenManager{
		secretKey:            secretKey,
		issuer:               issuer,
		audience:             audience,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (tm *TokenManager) GetRefreshTokenDuration() time.Duration {
	return tm.refreshTokenDuration
}

func (tm *TokenManager) GenerateAccessToken(userID, sessionID string, tokenVersion int) (string, error) {
	claims := &Claims{
		UserID:       userID,
		SessionID:    sessionID,
		TokenVersion: tokenVersion,
		TokenType:    "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    tm.issuer,
			Audience:  []string{tm.audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.secretKey))
}

func (tm *TokenManager) ValidateAccessToken(tokenStr string) (any, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(tm.secretKey), nil
		},
		jwt.WithIssuer(tm.issuer),
		jwt.WithAudience(tm.audience),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != "access" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random refresh token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func HashRefreshToken(refreshToken string) string {
	hash := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(hash[:])
}
