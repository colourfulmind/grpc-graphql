package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"ozon/internal/storage"
	"ozon/pkg/logger/sl"
	"time"

	"ozon/internal/domain/models"
)

type Resolver struct {
	DB  *sql.DB
	log *slog.Logger
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	SSLMode  string `yaml:"ssl_mode"`
}

const timeout = 30 * time.Second

func New(p Postgres, log *slog.Logger) (*MutationResolver, error) {
	const op = "internal/storage/graphql.New"

	conn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		p.Host, p.Port, p.User, p.DBName, p.Password, p.SSLMode)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	return &MutationResolver{
		&Resolver{
			DB:  db,
			log: log,
		},
	}, nil
}

func (r *Resolver) Mutation() MutationResolver {
	return MutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return QueryResolver{r}
}

type QueryResolver struct{ *Resolver }
type MutationResolver struct{ *Resolver }

// SaveUser grpcurl -plaintext -d '{"email": "test@test.com", "password": "1234567890"}' localhost:8080 sso.SSO/RegisterNewUser
func (r *MutationResolver) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "internal/storage/graphql.SaveUser"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var userID int64
	go func() {
		err := r.DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
		if errors.Is(err, sql.ErrNoRows) {
			err = r.DB.QueryRow("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id", email, string(passHash)).Scan(&userID)
			if err != nil {
				errQuery <- struct{}{}
			} else {
				done <- struct{}{}
			}
		} else {
			if !errors.Is(err, sql.ErrNoRows) {
				errExists <- struct{}{}
			} else {
				errQuery <- struct{}{}
			}
		}
		close(errQuery)
		close(errExists)
		close(done)
	}()

	select {
	case <-errExists:
		return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	case <-errQuery:
		return 0, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-ctx.Done():
		return 0, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return userID, nil
	}
}

func (r *MutationResolver) ProvideUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "internal/storage/graphql.ProvideUserByEmail"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var user models.User
	go func() {
		err := r.DB.QueryRow("SELECT (id, name, email) FROM users WHERE email = $1", email).Scan(&user)
		if err != nil {
			errQuery <- struct{}{}
		} else {
			if !errors.Is(err, sql.ErrNoRows) {
				errExists <- struct{}{}
			} else {
				done <- struct{}{}
			}
		}
		close(errQuery)
		close(errExists)
		close(done)
	}()

	select {
	case <-errExists:
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	case <-errQuery:
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return user, nil
	}
}

func (r *MutationResolver) ProvideUserByID(ctx context.Context, id int64) (models.User, error) {
	const op = "internal/storage/graphql.ProvideUserByEmail"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var user models.User
	go func() {
		err := r.DB.QueryRow("SELECT (id, name, email) FROM users WHERE id = $1", id).Scan(&user)
		if err != nil {
			errQuery <- struct{}{}
		} else {
			if !errors.Is(err, sql.ErrNoRows) {
				errExists <- struct{}{}
			} else {
				done <- struct{}{}
			}
		}
		close(errQuery)
		close(errExists)
		close(done)
	}()

	select {
	case <-errExists:
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	case <-errQuery:
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return user, nil
	}
}

func (r *MutationResolver) ProvideApp(ctx context.Context, appID int32) (models.App, error) {
	const op = "internal/storage/graphql.ProvideApp"
	panic("not implemented")
}

// SavePost grpcurl -plaintext -d '{"user_id": 1, "title": "test", "content": "Hello, world!", "allow_comments": true}' localhost:8080 posts.Posts/PostNew
func (r *MutationResolver) SavePost(ctx context.Context, userID int64, title, content string, allowComments bool) (int64, time.Time, error) {
	const op = "internal/storage/graphql.SavePost"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var post models.Post
	go func() {
		err := r.DB.QueryRow("SELECT id FROM posts WHERE user_id = $1 AND title = $2", userID, title).Scan(&post)
		if errors.Is(err, sql.ErrNoRows) {
			err = r.DB.QueryRow(
				"INSERT INTO posts (title, content, comments, user_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
				title, content, allowComments, userID).Scan(&post.ID, &post.CreatedAt)
			if err != nil {
				errQuery <- struct{}{}
			} else {
				done <- struct{}{}
			}

		} else {
			if !errors.Is(err, sql.ErrNoRows) {
				errExists <- struct{}{}
			} else {
				errQuery <- struct{}{}
			}
		}
		close(errQuery)
		close(errExists)
		close(done)
	}()

	select {
	case <-errExists:
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, storage.ErrPostExists)
	case <-errQuery:
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-ctx.Done():
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return post.ID, post.CreatedAt, nil
	}
}

// ProvidePost grpcurl -plaintext -d '{"id": 1}' localhost:8080 posts.Posts/GetPostByID
func (r *MutationResolver) ProvidePost(ctx context.Context, postID int64) (models.Post, error) {
	const op = "internal/storage/graphql.ProvidePost"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var post models.Post
	go func() {
		err := r.DB.QueryRow(
			"SELECT id, title, content, created_at, user_id, comments FROM posts WHERE id = $1", postID).Scan(
			&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UserID, &post.AllowComments)
		if errors.Is(err, sql.ErrNoRows) {
			errExists <- struct{}{}
		} else {
			if err != nil {
				errQuery <- struct{}{}
			} else {
				done <- struct{}{}
			}
		}
		close(errQuery)
		close(errExists)
		close(done)
	}()

	select {
	case <-errExists:
		return models.Post{}, fmt.Errorf("%s: %w", op, storage.ErrPostNotFound)
	case <-errQuery:
		return models.Post{}, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-ctx.Done():
		return models.Post{}, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return post, nil
	}
}

func (r *MutationResolver) ProvideAllPosts(ctx context.Context, page int64) ([]models.Post, error) {
	panic("not implemented")
}

// SaveComment grpcurl -plaintext -d '{"user_id": 1, "post_id": 1, "content": "test"}' localhost:8080 comments.Comments/PostNewComment
func (r *MutationResolver) SaveComment(ctx context.Context, userID, postID int64, content string) (int64, time.Time, error) {
	const op = "internal/storage/graphql.SaveComment"
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errNotAllowed := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var comment models.Comment
	var commentsAllowed bool
	go func() {
		err := r.DB.QueryRow("SELECT comments FROM posts WHERE id = $2 AND user_id = $1", postID, userID).Scan(&commentsAllowed)
		if errors.Is(err, sql.ErrNoRows) {
			errExists <- struct{}{}
		} else {
			if commentsAllowed {
				err = r.DB.QueryRow(
					"INSERT INTO comments (content, user_id, post_id, commemt_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
					content, userID, postID, 0).Scan(&comment.ID, &comment.CreatedAt)
				if err != nil {
					r.log.Info("error", sl.Err(err))
					errQuery <- struct{}{}
				} else {
					done <- struct{}{}
				}
			} else {
				errNotAllowed <- struct{}{}
			}
		}
		close(errQuery)
		close(errExists)
		close(errNotAllowed)
		close(done)
	}()

	select {
	case <-errExists:
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, storage.ErrPostNotFound)
	case <-errQuery:
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-errNotAllowed:
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, storage.ErrNotAllowed)
	case <-ctx.Done():
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return comment.ID, comment.CreatedAt, nil
	}
}

// SaveCommentToComment grpcurl -plaintext -d '{"user_id": 1, "post_id": 1, "comment_id": 1, "content": "test"}' localhost:8080 comments.Comments/PostCommentToComment
func (r *MutationResolver) SaveCommentToComment(ctx context.Context, userID, postID, commentID int64, content string) (int64, time.Time, error) {
	const op = "internal/storage/graphql.SaveCommentToComment"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var comment models.Comment
	go func() {
		err := r.DB.QueryRow(
			"INSERT INTO comments (content, user_id, post_id, commemt_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
			content, userID, postID, commentID).Scan(&comment.ID, &comment.CreatedAt)
		if err != nil {
			errQuery <- struct{}{}
		} else {
			done <- struct{}{}
		}
		close(errQuery)
		close(errExists)
		close(done)
	}()

	select {
	case <-errExists:
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, storage.ErrPostNotFound)
	case <-errQuery:
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-ctx.Done():
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return comment.ID, comment.CreatedAt, nil
	}
	panic("not implemented")
}

func (r *MutationResolver) ProvideComment(ctx context.Context, postID int64) ([]models.Comment, error) {
	const op = "internal/storage/graphql.ProvideComment"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	panic("not implemented")
}
