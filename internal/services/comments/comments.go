package comments

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"ozon/internal/domain/models"
	"ozon/internal/storage"
	"ozon/pkg/logger/sl"
)

// Comment - структура для работы с комментарием
type Comment struct {
	Log             *slog.Logger
	CommentSaver    CommentSaver
	CommentProvider CommentProvider
}

// CommentSaver - интерфейс для сохранения комментария в базе данных
type CommentSaver interface {
	SaveComment(ctx context.Context, userID, postID int64, content string) (int64, time.Time, error)
	SaveCommentToComment(ctx context.Context, userID, postID, commentID int64, content string) (int64, time.Time, error)
}

// CommentProvider - интерфейс для получения комментариев
type CommentProvider interface {
	ProvideComment(ctx context.Context, postID int64) ([]models.Comment, error)
}

var (
	ErrConnectionTime = errors.New("connection time expired")
	ErrGetComments    = errors.New("failed to get comments")
	ErrFoundComments  = errors.New("no comments found")
	ErrNotAllowed     = errors.New("not allowed")
)

// New возвращает структуру для работы с комментарием
func New(log *slog.Logger, commentSaver CommentSaver, commentProvider CommentProvider) *Comment {
	return &Comment{
		Log:             log,
		CommentSaver:    commentSaver,
		CommentProvider: commentProvider,
	}
}

// PostNewComment публикует новый комментарий к посту
// Если комментарий превышает максимальную длину, то возвращает ошибку
func (c *Comment) PostNewComment(ctx context.Context, userID, postID int64, content string) (int64, time.Time, error) {
	const op = "internal/services/comments/comments.PostNewComment"

	log := c.Log.With(slog.String("op", op), slog.Int64("user_id", userID))
	log.Info("attempting to create new comment")

	commentID, created, err := c.CommentSaver.SaveComment(ctx, userID, postID, content)
	if err != nil {
		if errors.Is(err, storage.ErrConnectionTime) {
			log.Error("connection time expired", sl.Err(err))
			return 0, time.Time{}, fmt.Errorf("%s: %w", op, ErrConnectionTime)
		}

		if errors.Is(err, storage.ErrNotAllowed) {
			log.Warn("comments not allowed", sl.Err(err))
			return 0, time.Time{}, fmt.Errorf("%s: %w", op, ErrNotAllowed)
		}

		log.Error("could not save comment", sl.Err(err))
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	return commentID, created, nil
}

// PostCommentToComment публикует новый комментарий к комментарию
// Если комментарий превышает максимальную длину, то возвращает ошибку
func (c *Comment) PostCommentToComment(ctx context.Context, userID, postID, commentID int64, content string) (int64, time.Time, error) {
	const op = "internal/services/comments/comments.PostCommentToComment"

	log := c.Log.With(slog.String("op", op), slog.Int64("comment_id", commentID))
	log.Info("attempting to comment")

	newCommentID, created, err := c.CommentSaver.SaveCommentToComment(ctx, userID, postID, commentID, content)
	if err != nil {
		if errors.Is(err, storage.ErrConnectionTime) {
			log.Error("connection time expired", sl.Err(err))
			return 0, time.Time{}, fmt.Errorf("%s: %w", op, ErrConnectionTime)
		}

		log.Error("could not save comment", sl.Err(err))
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	return newCommentID, created, nil
}

// GetComments возвращает все комментарии к посту
// Возвращает ошибку, если не удалось получить комментарии из БД
func (c *Comment) GetComments(ctx context.Context, postID int64) ([]models.Comment, error) {
	const op = "internal/services/comments/comments.GetComments"

	log := c.Log.With(slog.String("op", op), slog.Int64("post_id", postID))
	log.Info("attempting to get all comments")

	comments, err := c.CommentProvider.ProvideComment(ctx, postID)
	if err != nil {
		if errors.Is(err, storage.ErrGetComments) {
			log.Error("could not get comments", sl.Err(err))
			return []models.Comment{}, fmt.Errorf("%s: %w", op, ErrGetComments)
		}

		if errors.Is(err, storage.ErrFoundComments) {
			log.Info("could not find comments", sl.Err(err))
			return []models.Comment{}, fmt.Errorf("%s: %w", op, ErrFoundComments)
		}

		if errors.Is(err, storage.ErrConnectionTime) {
			log.Error("connection time expired", sl.Err(err))
			return []models.Comment{}, fmt.Errorf("%s: %w", op, ErrConnectionTime)
		}

		log.Error("could not get comments", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}
