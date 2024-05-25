package sso

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"regexp"

	"ozon/internal/ewrap"
	"ozon/internal/services/sso"
	server "ozon/protos/gen/go/sso"
)

// ServerSSO наследует интерфейс из protobuf и включает собственную реализацию
type ServerSSO struct {
	server.UnimplementedSSOServer
	sso SSO
}

// SSO описывает методы для сервиса авторизации
type SSO interface {
	RegisterNewUser(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string, appID int32) (string, error)
}

// Register регистрирует sso сервер
func Register(s *grpc.Server, sso SSO) {
	server.RegisterSSOServer(s, &ServerSSO{
		sso: sso,
	})
}

// RegisterNewUser регистрирует нового пользователя
// Возвращает ошибку, если пользователь уже существует, введены неправильные логин или пароль
// либо при возникновении ошибки при подключении к базе данных
func (s *ServerSSO) RegisterNewUser(ctx context.Context, req *server.RegisterRequest) (*server.RegisterResponse, error) {
	const op = "internal/grpc/sso.RegisterNewUser"

	if err := ValidateCredentials(req.GetEmail(), req.GetPassword()); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := s.sso.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, sso.ErrUserExists) {
			return nil, ewrap.UserAlreadyExists
		}

		if errors.Is(err, sso.ErrConnectionTime) {
			return nil, ewrap.ErrConnectionTime
		}

		return nil, ewrap.InternalError
	}

	return &server.RegisterResponse{
		UserId: userID,
	}, nil
}

// Login авторизовывает пользователя
// Возвращает ошибку, если введены неправильные логин или пароль
// либо при возникновении ошибки при подключении к базе данных
func (s *ServerSSO) Login(ctx context.Context, req *server.LoginRequest) (*server.LoginResponse, error) {
	const op = "internal/grpc/sso.Login"

	if err := ValidateCredentials(req.GetEmail(), req.GetPassword()); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := s.sso.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, sso.ErrInvalidCredentials) {
			return nil, ewrap.ErrInvalidCredentials
		}

		if errors.Is(err, sso.ErrConnectionTime) {
			return nil, ewrap.ErrConnectionTime
		}

		return nil, ewrap.InternalError
	}

	return &server.LoginResponse{
		Token: token,
	}, nil
}

// ValidateCredentials проверяет формат логина и пароля
func ValidateCredentials(email, password string) error {
	matched, err := regexp.MatchString(`([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)`, email)
	if err != nil {
		return ewrap.ErrParsingRegex
	}

	if !matched || email == "" {
		return ewrap.ErrInvalidEmail
	}

	if password == "" || len(password) < 8 {
		return ewrap.ErrPasswordRequired
	}

	return nil
}
