package store

import (
	"context"

	"github.com/yusuf-cirak/social/internal/db"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}
type UserStore struct {
	db *db.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := s.db.QueryRow(ctx, query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
	user := &User{}
	err := s.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
