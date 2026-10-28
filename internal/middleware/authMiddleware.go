package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aryansaud-80/rosterpulse/internal/shared/httpx"
)

type TokenManager interface {
	ValidateAccessToken(tokenStr string) (any, error)
}

type jwtContextKey string

const claimskey jwtContextKey = "jwtClaims"

func AuthMiddleware(tm TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				httpx.Error(w, http.StatusUnauthorized, "invalid authorization format", httpx.ErrCodeUnauthorized)
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := tm.ValidateAccessToken(token)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "invalid or expired token", httpx.ErrCodeUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimskey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) (any, bool) {
	claims := ctx.Value(claimskey)
	if claims == nil {
		return nil, false
	}
	return claims, true
}
