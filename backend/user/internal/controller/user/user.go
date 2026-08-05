package user

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/mail"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const maxRequestBodySize = 1 << 20

//go:generate mockgen -source=user.go -destination=mocks/user_mocks.go -package=mocks

type authService interface {
	Register(ctx context.Context, user entity.User) (*entity.JWT, error)
	Login(ctx context.Context, email, password string) (*entity.JWT, error)
	Refresh(ctx context.Context, refreshToken string) (*entity.JWT, error)
}

type userService interface {
	Profile(ctx context.Context, userID string) (*entity.User, error)
	GetUsers(ctx context.Context, userIDs []string) ([]entity.User, error)
	UpdateCoins(ctx context.Context, userID string, coins int64) (*entity.User, error)
}

type controller struct {
	logger      *zap.Logger
	authService authService
	userService userService
}

func NewController(logger *zap.Logger, authService authService, userService userService) *controller {
	return &controller{
		logger:      logger,
		authService: authService,
		userService: userService,
	}
}

func (c *controller) InitRoutes(router *mux.Router, auth func(http.Handler) http.Handler) {
	api := router.PathPrefix("/api/v1").Subrouter()
	internal := router.PathPrefix("/internal").Subrouter()

	authRouter := api.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", c.Register).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", c.Login).Methods(http.MethodPost)
	authRouter.HandleFunc("/refresh", c.Refresh).Methods(http.MethodPost)

	profile := api.PathPrefix("/profile").Subrouter()
	profile.Use(auth)
	profile.HandleFunc("", c.Profile).Methods(http.MethodGet)

	internal.HandleFunc("/usernames", c.Usernames).Methods(http.MethodPost)
	internal.HandleFunc("/update-coins", c.UpdateCoins).Methods(http.MethodPut)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain exactly one JSON object")
	}
	return nil
}

func isValidUsername(username string) bool {
	length := len([]rune(username))
	return length >= 2 && length <= 40
}

func isValidEmail(email string) bool {
	parsed, err := mail.ParseAddress(email)
	return err == nil && parsed.Address == email
}

func isValidPassword(password string) bool {
	length := len([]byte(password))
	return length >= 8 && length <= 72
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type authResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type profileResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Coins    uint64 `json:"coins"`
}

type usernamesRequest struct {
	UserIDs []string `json:"user_ids"`
}

type usernameResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type usernamesResponse struct {
	Users []usernameResponse `json:"users"`
}

type updateCoinsRequest struct {
	UserID string `json:"user_id"`
	Delta  int64  `json:"delta"`
}

type updateCoinsResponse struct {
	UserID string `json:"user_id"`
	Coins  uint64 `json:"coins"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errCode string

var (
	errUserAlreadyExists   errCode = "USER_ALREADY_EXISTS"
	errValidationError     errCode = "VALIDATION_ERROR"
	errInternalError       errCode = "INTERNAL_ERROR"
	errInvalidCredentials  errCode = "INVALID_CREDENTIALS"
	errInvalidRefreshToken errCode = "INVALID_REFRESH_TOKEN"
	errUnauthorized        errCode = "UNAUTHORIZED"
	errUserNotFound        errCode = "USER_NOT_FOUND"
	errInsufficientCoins   errCode = "INSUFFICIENT_COINS"
)

func (e errCode) message() string {
	switch e {
	case errUserAlreadyExists:
		return "Пользователь с такими данными уже существует"
	case errValidationError:
		return "Проверьте переданные данные"
	case errInternalError:
		return "Внутренняя ошибка сервера"
	case errInvalidCredentials:
		return "Неверный email или пароль"
	case errInvalidRefreshToken:
		return "Сессия истекла. Выполните вход снова"
	case errUnauthorized:
		return "Требуется повторная авторизация"
	case errUserNotFound:
		return "Пользователь не найден"
	case errInsufficientCoins:
		return "Недостаточно монет"

	default:
		return "Unknown error"
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code errCode) {
	writeJSON(w, status, errorResponse{Code: string(code), Message: code.message()})
}
