package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/artlink52/ecommerce_backend/pkg/logger"
	"github.com/artlink52/ecommerce_backend/services/api-gateway/internal/transport/response"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrHeaderMissing = errors.New("authorization header missing")
	ErrHeaderFormat  = errors.New("authorization header format error")
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidClaims = errors.New("invalid claims")
)

type ContextKey string

const UserIDKey ContextKey = "user_id"

type Middleware func(http.Handler) http.Handler

func AuthMiddleware(jwtSecret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			responseHandler := response.NewHTTPResponseHandler(log, w)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Error("Authorization header missing")
				responseHandler.ErrorResponse(ErrHeaderMissing, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Error("Authorization header format error", slog.String("authHeader", authHeader))
				responseHandler.ErrorResponse(ErrHeaderFormat, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				responseHandler.ErrorResponse(err, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Error("Authorization header claims invalid", slog.String("claims", fmt.Sprintf("%v", claims)))
				responseHandler.ErrorResponse(ErrInvalidClaims, "invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(int64)
			if !ok {
				log.Error("Authorization header claims invalid", slog.String("claims", fmt.Sprintf("%v", claims)))
				responseHandler.ErrorResponse(ErrInvalidClaims, "invalid token claims", http.StatusUnauthorized)
				return
			}

			newCtx := context.WithValue(ctx, UserIDKey, int64(userID))
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}
