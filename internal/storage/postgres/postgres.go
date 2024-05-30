package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"time"

	"ozon/internal/domain/models"
	"ozon/internal/storage"
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
	const op = "internal/storage/postgres.New"

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

type QueryResolver struct{ *Resolver }
type MutationResolver struct{ *Resolver }

func (r *MutationResolver) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "internal/storage/postgres.SaveUser"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var userID int64
	go func() {
		err := r.DB.QueryRowContext(ctx,
			"SELECT id FROM users WHERE email = $1", email).Scan(&userID)
		if errors.Is(err, sql.ErrNoRows) {
			err = r.DB.QueryRowContext(ctx,
				"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id", email, string(passHash)).Scan(&userID)
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
	const op = "internal/storage/postgres.ProvideUserByEmail"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var user models.User
	go func() {
		err := r.DB.QueryRowContext(ctx,
			"SELECT id, email, password FROM users WHERE email = $1", email).Scan(
			&user.ID, &user.Email, &user.Password)
		if errors.Is(err, sql.ErrNoRows) {
			errExists <- struct{}{}
		} else {
			if !errors.Is(err, sql.ErrNoRows) {
				done <- struct{}{}
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
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	case <-errQuery:
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return user, nil
	}
}

func (r *MutationResolver) ProvideUserByID(ctx context.Context, id int64) (models.User, error) {
	const op = "internal/storage/postgres.ProvideUserByEmail"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var user models.User
	go func() {
		err := r.DB.QueryRowContext(ctx,
			"SELECT (id, name, email) FROM users WHERE id = $1", id).Scan(&user)
		if errors.Is(err, sql.ErrNoRows) {
			errExists <- struct{}{}
		} else {
			if !errors.Is(err, sql.ErrNoRows) {
				done <- struct{}{}
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
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	case <-errQuery:
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return user, nil
	}
}

func (r *MutationResolver) SavePost(ctx context.Context, userID int64, title, content string, allowComments bool) (int64, time.Time, error) {
	const op = "internal/storage/postgres.SavePost"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	r.log.Info("data", slog.Int64("user_id", userID), slog.String("title", title),
		slog.String("content", content), slog.Bool("allowComments", allowComments))

	var post models.Post
	go func() {
		err := r.DB.QueryRowContext(ctx,
			"SELECT id FROM posts WHERE user_id = $1 AND title = $2", userID, title).Scan(&post)
		if errors.Is(err, sql.ErrNoRows) {
			err = r.DB.QueryRow(
				"INSERT INTO posts (title, content, comments_allowed, user_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
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

func (r *MutationResolver) ProvidePost(ctx context.Context, postID int64) (models.Post, error) {
	const op = "internal/storage/postgres.ProvidePost"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var post models.Post
	go func() {
		err := r.DB.QueryRowContext(ctx,
			"SELECT id, title, content, comments_allowed, created_at, user_id FROM posts WHERE id = $1", postID).Scan(
			&post.ID, &post.Title, &post.Content, &post.AllowComments, &post.CreatedAt, &post.UserID)
		r.log.Info("posts", slog.Any("posts", post))
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
	const op = "internal/storage/postgres.ProvideAllPosts"
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var posts []models.Post
	go func() {
		rows, err := r.DB.QueryContext(ctx,
			"SELECT id, title, content, comments_allowed, created_at, user_id  FROM posts ORDER BY id LIMIT 10 OFFSET $1;", (page-1)*10)
		defer rows.Close()

		if err != nil {
			errQuery <- struct{}{}
		} else {
			for rows.Next() {
				var row models.Post

				err = rows.Scan(&row.ID, &row.Title, &row.Content, &row.AllowComments, &row.CreatedAt, &row.UserID)
				if err != nil {
					errQuery <- struct{}{}
				}
				posts = append(posts, row)
			}
			if err = rows.Err(); err != nil {
				errQuery <- struct{}{}
			}

			done <- struct{}{}
		}
		close(errQuery)
		close(done)
	}()

	select {
	case <-errQuery:
		return []models.Post{}, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-ctx.Done():
		return []models.Post{}, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return posts, nil
	}
}

func (r *MutationResolver) SaveComment(ctx context.Context, userID, postID int64, content string) (int64, time.Time, error) {
	const op = "internal/storage/postgres.SaveComment"
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errNotAllowed := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var comment models.Comment
	var commentsAllowed bool
	go func() {
		err := r.DB.QueryRow("SELECT comments_allowed FROM posts WHERE id = $1", postID).Scan(&commentsAllowed)
		if errors.Is(err, sql.ErrNoRows) {
			errExists <- struct{}{}
		} else {
			if commentsAllowed {
				err = r.DB.QueryRow(
					"INSERT INTO comments (content, user_id, post_id) VALUES ($1, $2, $3) RETURNING id, created_at",
					content, userID, postID).Scan(&comment.ID, &comment.CreatedAt)
				if err != nil {
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

func (r *MutationResolver) SaveCommentToComment(ctx context.Context, userID, postID, parentID int64, content string) (int64, time.Time, error) {
	const op = "internal/storage/postgres.SaveCommentToComment"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var comment models.Comment
	go func() {
		err := r.DB.QueryRow(
			"INSERT INTO comments (content, user_id, post_id, parent_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
			content, userID, postID, parentID).Scan(&comment.ID, &comment.CreatedAt)
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
}

func (r *MutationResolver) ProvideComment(ctx context.Context, postID, parentID int64) ([]models.Comment, error) {
	const op = "internal/storage/postgres.ProvideComment"

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{}, 1)
	errExists := make(chan struct{}, 1)
	errQuery := make(chan struct{}, 1)

	var comments []models.Comment

	go func() {
		var query string
		var rows *sql.Rows
		var err error
		if parentID == 0 {
			query = "SELECT id, content, user_id, post_id, created_at FROM comments WHERE post_id = $1 AND parent_id IS NULL;"
			rows, err = r.DB.QueryContext(ctx, query, postID)
		} else {
			query = "SELECT id, content, user_id, post_id, parent_id, created_at FROM comments WHERE post_id = $1 AND parent_id = $2;"
			rows, err = r.DB.QueryContext(ctx, query, postID, parentID)
		}

		defer rows.Close()
		if errors.Is(err, sql.ErrNoRows) {
			errExists <- struct{}{}
		} else {
			if err != nil {
				errQuery <- struct{}{}
			} else {
				for rows.Next() {
					var row models.Comment
					if parentID == 0 {
						err = rows.Scan(&row.ID, &row.Content, &row.UserID, &row.PostID, &row.CreatedAt)
					} else {
						err = rows.Scan(&row.ID, &row.Content, &row.UserID, &row.PostID, &row.ParentID, &row.CreatedAt)
					}
					if err != nil {
						errQuery <- struct{}{}
					}
					comments = append(comments, row)
				}
				if err = rows.Err(); err != nil {
					errQuery <- struct{}{}
				}

				done <- struct{}{}
			}
		}
		close(errQuery)
		close(errExists)
		close(done)
	}()

	select {
	case <-errExists:
		return []models.Comment{}, fmt.Errorf("%s: %w", op, storage.ErrFoundComments)
	case <-errQuery:
		return []models.Comment{}, fmt.Errorf("%s: %w", op, storage.ErrQuery)
	case <-ctx.Done():
		return []models.Comment{}, fmt.Errorf("%s: %w", op, storage.ErrConnectionTime)
	case <-done:
		return comments, nil
	}
}
