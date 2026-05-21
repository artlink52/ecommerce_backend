package jwt

import (
	"time"

	"github.com/artlink52/ecommerce_backend/services/user-service/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	return token.SignedString([]byte(secret))
}
