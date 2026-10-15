package middleware

import (
	"TaskManagerApi/auth"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(
				w,
				"authorization required",
				http.StatusUnauthorized,
			)
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(
				w,
				"invalid authorization header",
				http.StatusUnauthorized,
			)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return auth.GetJWTSecret(), nil
			},
		)

		if err != nil || !token.Valid {
			http.Error(
				w,
				"invalid token",
				http.StatusUnauthorized,
			)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(
				w,
				"invalid token claims",
				http.StatusUnauthorized,
			)
			return
		}

		userID, ok := claims["user_id"]

		if !ok {
			http.Error(
				w,
				"user id not found",
				http.StatusUnauthorized,
			)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			"user_id",
			userID,
		)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
