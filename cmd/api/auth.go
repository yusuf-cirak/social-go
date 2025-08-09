package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/yusuf-cirak/social/internal/store"
)

type authConfig struct {
	Secret         string
	Issuer         string
	Audience       string
	AccessTokenTTL time.Duration
}

type currentUserKey string

const currentUserCtxKey currentUserKey = "current_user"

// authMiddleware validates the bearer token and loads the current user into context.
func (app *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			writeJSONError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeJSONError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		token := parts[1]
		claims, err := app.jwt.ParseAndValidate(token)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Optionally load user from DB to ensure it still exists
		user, err := app.store.Users.GetByID(r.Context(), claims.UserID)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "user not found")
			return
		}

		ctx := context.WithValue(r.Context(), currentUserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCurrentUser(ctx context.Context) *store.User {
	user, ok := ctx.Value(currentUserCtxKey).(*store.User)
	if !ok {
		return nil
	}
	return user
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginHandler authenticates a user and returns a JWT access token.
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var payload loginPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		app.badRequest(w, r, errors.New("email and password are required"))
		return
	}

	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		app.notFound(w, r, err)
		return
	}

	// NOTE: In production, compare password hashes. This demo uses plaintext from seed data.
	if user.Password != payload.Password {
		writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := app.jwt.GenerateToken(user.ID, user.Username, app.config.auth.AccessTokenTTL)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	_ = app.jsonResponse(w, http.StatusOK, map[string]string{"token": token})
}
