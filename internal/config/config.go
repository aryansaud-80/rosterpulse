package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Env          string
	DATABASE_URL string
	JWT          JWTConfig
}

type JWTConfig struct {
	SecretKey            string
	Issuer               string
	Audience             string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func MustLoad() Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required")
	}

	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is required")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		panic("DATABASE_URL is required")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		panic("JWT_SECRET_KEY is required")
	}

	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		panic("JWT_ISSUER is required")
	}

	jwtAudience := os.Getenv("JWT_AUDIENCE")
	if jwtAudience == "" {
		panic("JWT_AUDIENCE is required")
	}

	jwtAccessTokenDurationStr := os.Getenv("JWT_ACCESS_TOKEN_DURATION")
	if jwtAccessTokenDurationStr == "" {
		panic("JWT_ACCESS_TOKEN_DURATION is required")
	}
	jwtAccessTokenDuration, err := time.ParseDuration(jwtAccessTokenDurationStr)
	if err != nil {
		panic("JWT_ACCESS_TOKEN_DURATION is invalid: " + err.Error())
	}

	jwtRefreshTokenDurationStr := os.Getenv("JWT_REFRESH_TOKEN_DURATION")
	if jwtRefreshTokenDurationStr == "" {
		panic("JWT_REFRESH_TOKEN_DURATION is required")
	}
	jwtRefreshTokenDuration, err := time.ParseDuration(jwtRefreshTokenDurationStr)
	if err != nil {
		panic("JWT_REFRESH_TOKEN_DURATION is invalid: " + err.Error())
	}

	jwt := JWTConfig{
		SecretKey:            jwtSecretKey,
		Issuer:               jwtIssuer,
		Audience:             jwtAudience,
		AccessTokenDuration:  jwtAccessTokenDuration,
		RefreshTokenDuration: jwtRefreshTokenDuration,
	}

	return Config{
		Port:         port,
		Env:          env,
		DATABASE_URL: databaseURL,
		JWT:          jwt,
	}
}
