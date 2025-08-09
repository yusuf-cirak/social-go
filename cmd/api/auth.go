package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	iauth "github.com/yusuf-cirak/social/internal/auth"
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

// authorize returns a middleware that enforces a policy action on a resource extracted from the request.
func (app *application) authorize(action string, extract func(*http.Request) (iauth.Resource, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			current := getCurrentUser(r.Context())
			if current == nil {
				writeJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			res, err := extract(r)
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, "invalid resource")
				return
			}

			sub := iauth.Subject{UserID: current.ID}
			if app.policy.Authorize(sub, action, res) {
				next.ServeHTTP(w, r)
				return
			}

			writeJSONError(w, http.StatusForbidden, "forbidden")
		})
	}
}

// Resource helpers
func (app *application) resourcePostCreate(r *http.Request) (iauth.Resource, error) {
	return iauth.Resource{Type: "post"}, nil
}

func (app *application) resourcePostFromCtx(r *http.Request) (iauth.Resource, error) {
	p := getPostFromCtx(r)
	if p == nil {
		return iauth.Resource{}, errors.New("post not in context")
	}
	return iauth.Resource{Type: "post", OwnerID: p.UserID}, nil
}

func (app *application) resourceUserFromCtx(r *http.Request) (iauth.Resource, error) {
	u := getUserFromContext(r.Context())
	if u == nil {
		return iauth.Resource{}, errors.New("user not in context")
	}
	return iauth.Resource{Type: "user", OwnerID: u.ID}, nil
}
