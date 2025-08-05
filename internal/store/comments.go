package store

import (
	"context"
	"time"

	"github.com/yusuf-cirak/social/internal/db"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"`
}

type CommentStore struct {
	db *db.DB
}

func NewCommentStore(db *db.DB) *CommentStore {
	return &CommentStore{db: db}
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) 
	RETURNING id, created_at`
	err := s.db.QueryRow(ctx, query, comment.PostID, comment.UserID, comment.Content).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.id FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE post_id = $1
	ORDER BY c.created_at DESC
	`
	rows, err := s.db.Query(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		comment.User = User{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.User.Username, &comment.User.ID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (s *CommentStore) Delete(ctx context.Context, postID int64) error {
	query := `
	DELETE FROM comments
	WHERE post_id = $1
	`
	res, err := s.db.Exec(ctx, query, postID)

	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
