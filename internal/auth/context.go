package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/aryansaud-80/rosterpulse/internal/middleware"
)

func UserIDFromContext(ctx context.Context) (string, error) {
	rawClaims, ok := middleware.ClaimsFromContext(ctx)
	if !ok {
		return "", errors.New("unauthorized")
	}

	claims, ok := rawClaims.(*Claims)
	if !ok || claims == nil {
		return "", errors.New("invalid token claims")
	}

	userId := strings.TrimSpace(claims.UserID)
	if userId == "" {
		return "", errors.New("missing user id in claims")
	}

	return userId, nil
}
