package seed

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/yusuf-cirak/social/internal/store"
)

func Seed(store store.Storage) error {
	ctx := context.Background()

	users := generateUsers(100)

	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			return err
		}
	}

	posts := generatePosts(100, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			return err
		}
	}

	comments := generateComments(500, users, posts)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			return err
		}
	}

	return nil
}

func generateUsers(count int) []*store.User {
	users := make([]*store.User, count)
	now := time.Now().Format(time.RFC3339)

	for i := 0; i < count; i++ {
		var username, email, password strings.Builder

		username.WriteString("user")
		username.WriteString(strconv.Itoa(i))

		email.WriteString("user")
		email.WriteString(strconv.Itoa(i))
		email.WriteString("@example.com")

		password.WriteString("password")
		password.WriteString(strconv.Itoa(i))

		users[i] = &store.User{
			Username:  username.String(),
			Email:     email.String(),
			Password:  password.String(),
			CreatedAt: now,
		}
	}
	return users
}

func generatePosts(count int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, count)
	now := time.Now().Format(time.RFC3339)

	for i := 0; i < count; i++ {
		var content, title strings.Builder

		content.WriteString("This is post content ")
		content.WriteString(strconv.Itoa(i))

		title.WriteString("Post Title ")
		title.WriteString(strconv.Itoa(i))

		posts[i] = &store.Post{
			Content:   content.String(),
			Title:     title.String(),
			UserID:    users[i%len(users)].ID,
			Tags:      []string{"tag1", "tag2"},
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return posts
}

func generateComments(count int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, count)
	now := time.Now().UTC()

	for i := 0; i < count; i++ {
		var content strings.Builder
		content.WriteString("This is comment content ")
		content.WriteString(strconv.Itoa(i))

		comments[i] = &store.Comment{
			PostID:    posts[i%len(posts)].ID,
			UserID:    users[i%len(users)].ID,
			Content:   content.String(),
			CreatedAt: now,
		}
	}
	return comments
}
