package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/diegorezm/ticketing/internal/jwt"
	"github.com/diegorezm/ticketing/internal/responses"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var tokenStr string

		// 1️⃣ Try cookie
		if c, err := r.Cookie("access_token"); err == nil {
			tokenStr = c.Value
		}

		// 2️⃣ Fallback to Authorization header
		if tokenStr == "" {
			h := r.Header.Get("Authorization")
			if after, ok := strings.CutPrefix(h, "Bearer "); ok {
				tokenStr = after
			}
		}

		if tokenStr == "" {
			responses.Fail(w, http.StatusUnauthorized, "missing token")
			return
		}

		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			responses.Fail(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), jwt.UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
