package store

import (
	"context"

	"github.com/yusuf-cirak/social/internal/db"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *db.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `
	INSERT INTO followers (user_id, follower_id)
	VALUES ($1, $2)`
	_, err := s.db.Exec(ctx, query, userID, followerID)
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID int64, userID int64) error {
	query := `
	DELETE FROM followers
	WHERE user_id = $1 AND follower_id = $2`
	_, err := s.db.Exec(ctx, query, userID, followerID)
	return err
}
