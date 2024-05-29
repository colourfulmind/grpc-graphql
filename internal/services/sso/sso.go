package sso

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"

	"ozon/internal/domain/models"
	"ozon/internal/storage"
	"ozon/pkg/jwt"
	"ozon/pkg/logger/sl"
)

var (
	ErrUserExists         = errors.New("user does not exists")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrConnectionTime     = errors.New("cannot connect to database")
)

// SSO - структура для работы с пользователем
type SSO struct {
	Log          *slog.Logger
	UserSaver    UserSaver
	UserProvider UserProvider
	TokenTTl     time.Duration
}

// UserSaver - интерфейс для сохранения пользователя в базе данных
type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

// UserProvider - интерфейс для получения информации о пользователе
type UserProvider interface {
	ProvideUserByEmail(ctx context.Context, email string) (models.User, error)
	ProvideUserByID(ctx context.Context, id int64) (models.User, error)
}

// New возвращает структуру для работы с пользователем
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, tokenTTL time.Duration) *SSO {
	return &SSO{
		Log:          log,
		UserSaver:    userSaver,
		UserProvider: userProvider,
		TokenTTl:     tokenTTL,
	}
}

// RegisterNewUser генерирует хэш пароля и сохраняет пользователя в базу данных
// Если пользователь уже существует - возвращает ошибку
func (sso *SSO) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "internal/services/sso.RegisterNewUser"

	log := sso.Log.With(slog.String("op", op), slog.String("email", email))
	log.Info("attempting to register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	userID, err := sso.UserSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		if errors.Is(err, storage.ErrConnectionTime) {
			log.Error("connection time expired", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrConnectionTime)
		}

		log.Error("failed to register new user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user is registered", slog.Int64("user_id", userID))
	return userID, nil
}

// Login ищет в базе данных пользователя с отправленным в запросе email-адресом
// генерирует хэш пароля и сравнивает его с хранящимся в базе данных
// Возвращает ошибку, если пользовтель не найден или передан неверный пароль
func (sso *SSO) Login(ctx context.Context, email, password string) (string, error) {
	const op = "internal/services/sso.Login"

	log := sso.Log.With(slog.String("op", op), slog.String("email", email))
	log.Info("attempting to login user")

	user, err := sso.UserProvider.ProvideUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		if errors.Is(err, storage.ErrConnectionTime) {
			log.Error("connection time expired", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrConnectionTime)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.New(user, sso.TokenTTl)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user successfully logged in", slog.String("email", email))

	return token, nil
}
