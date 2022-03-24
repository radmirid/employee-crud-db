package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/radmirid/employee-crud-db/internal/domain"
	logger "github.com/radmirid/grpc-logger/pkg/domain"
	"github.com/sirupsen/logrus"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepo interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

type SessionsRepository interface {
	Create(ctx context.Context, token domain.RefreshSession) error
	Get(ctx context.Context, token string) (domain.RefreshSession, error)
}

type LoggerClient interface {
	SendLogRequest(ctx context.Context, req logger.LogItem) error
}

type Users struct {
	repo         UsersRepo
	sessionsRepo SessionsRepository
	hasher       PasswordHasher
	loggerClient LoggerClient
	hmacSecret   []byte
}

func NewUsers(repo UsersRepo, sessionsRepo SessionsRepository, loggerClient LoggerClient, hasher PasswordHasher, secret []byte) *Users {
	return &Users{
		repo:         repo,
		sessionsRepo: sessionsRepo,
		hasher:       hasher,
		loggerClient: loggerClient,
		hmacSecret:   secret,
	}
}

func (u *Users) SignUp(ctx context.Context, inp domain.SignUpInput) error {
	password, err := u.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return err
	}

	user, err = u.repo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		return err
	}

	if err := u.loggerClient.SendLogRequest(ctx, logger.LogItem{
		Action:    logger.ACTION_REGISTER,
		Entity:    logger.ENTITY_USER,
		EntityID:  user.ID,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Users.SignUp",
		}).Error("failed to send log request:", err)
	}

	return nil
}

func (u *Users) SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error) {
	password, err := u.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := u.repo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		return "", "", err
	}

	return u.generateTokens(ctx, user.ID)

}

func (u *Users) ParseToken(ctx context.Context, token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		return u.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}

func (u *Users) generateTokens(ctx context.Context, userId int64) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userId)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	})

	accessToken, err := token.SignedString(u.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := u.sessionsRepo.Create(ctx, domain.RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (u *Users) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := u.sessionsRepo.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrorRefreshTokenExpired
	}

	return u.generateTokens(ctx, session.UserID)
}
