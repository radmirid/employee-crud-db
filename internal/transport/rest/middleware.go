package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type CtxValue int

const (
	ctxUserID CtxValue = iota
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Info()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authorizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			logError("authorizer", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := h.usersService.ParseToken(r.Context(), token)
		if err != nil {
			logError("authorizer", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserID, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerIntoParts := strings.Split(header, " ")
	if len(headerIntoParts) != 2 || headerIntoParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerIntoParts[1]) == 0 {
		return "", errors.New("empty token")
	}

	return headerIntoParts[1], nil
}
