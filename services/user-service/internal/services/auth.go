package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/artlink52/ecommerce_backend/pkg/logger"
	domainerrors "github.com/artlink52/ecommerce_backend/services/user-service/internal/domain/errors"
	"github.com/artlink52/ecommerce_backend/services/user-service/internal/domain/models"
	"github.com/artlink52/ecommerce_backend/services/user-service/lib/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error)
	User(ctx context.Context, email string) (models.User, error)
}

type AuthService struct {
	log            *slog.Logger
	tokenTTL       time.Duration
	jwtSecret      string
	userRepository UserRepository
}

func NewAuthService(
	log *slog.Logger,
	tokenTTL time.Duration,
	jwtSecret string,
	userRepository UserRepository,
) *AuthService {
	return &AuthService{
		log:            log,
		tokenTTL:       tokenTTL,
		jwtSecret:      jwtSecret,
		userRepository: userRepository,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.Register"
	log := s.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("registering user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := s.userRepository.SaveUser(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, domainerrors.ErrUserExists) {
			log.Warn("user already exists", logger.Err(err))
			return 0, fmt.Errorf("%s: %w", op, domainerrors.ErrUserExists)
		}
		log.Error("failed to save user", logger.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")
	return userID, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	const op = "auth.Login"
	log := s.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("fetching user")

	user, err := s.userRepository.User(ctx, email)
	if err != nil {
		if errors.Is(err, domainerrors.ErrUserNotFound) {
			log.Warn("user not found", logger.Err(err))
			return "", fmt.Errorf("%s: %w", op, domainerrors.ErrInvalidCredentials)
		}
		log.Error("failed to fetch user", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Warn("invalid credentials", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, domainerrors.ErrInvalidCredentials)
	}

	log.Info("successfully logged in")

	token, err := jwt.NewToken(user, s.tokenTTL, s.jwtSecret)
	if err != nil {
		log.Error("failed to generate token", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
