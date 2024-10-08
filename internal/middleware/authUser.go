package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/NeGat1FF/todolist-api/internal/models"
	"github.com/NeGat1FF/todolist-api/internal/utils"
)

func AuthUserMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), "Bearer", "", -1))
		if tokenString == "" {
			http.Error(rw, "no authorization token", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(rw, "failed to validate token", http.StatusUnauthorized)
			return
		}

		var user_id float64

		expTime, ok := claims["exp"].(float64)
		if !ok || time.Now().After(time.Unix(int64(expTime), 0)) {
			http.Error(rw, "token expired", http.StatusUnauthorized)
			return
		}
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			http.Error(rw, "invalid type of token", http.StatusUnauthorized)
			return
		}
		user_id, ok = claims["uid"].(float64)
		if !ok {
			http.Error(rw, "token with invalid claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserIDKey{}, int(user_id))
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	}
}
