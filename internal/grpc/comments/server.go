package comments

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"ozon/internal/domain/models"
	"ozon/internal/ewrap"
	"ozon/internal/services/comments"
	server "ozon/protos/gen/go/comments"
)

// ServerComments наследует интерфейс из protobuf и включает собственную реализацию
type ServerComments struct {
	server.UnimplementedCommentsServer
	comments Comments
}

// Comments - описыает методы для сервиса публикации комментариев
type Comments interface {
	PostNewComment(ctx context.Context, token string, postID int64, content string) (int64, time.Time, error)
	PostCommentToComment(ctx context.Context, token string, postID, commentID int64, content string) (int64, time.Time, error)
	GetComments(ctx context.Context, postID, parentID int64) ([]models.Comment, error)
}

// Register регистрирует sso сервер
func Register(s *grpc.Server, comments Comments) {
	server.RegisterCommentsServer(s, &ServerComments{
		comments: comments,
	})
}

// PostNewComment публикует новый комментарий к посту
func (s *ServerComments) PostNewComment(ctx context.Context, req *server.NewCommentRequest) (*server.NewCommentResponse, error) {
	const op = "internal/services/comments.PostNewComment"

	if err := ValidateContent(req.Content); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	commentID, createdAt, err := s.comments.PostNewComment(ctx, req.GetToken(), req.GetPostId(), req.GetContent())
	if err != nil {
		if errors.Is(err, comments.ErrConnectionTime) {
			return nil, ewrap.ErrConnectionTime
		}

		if errors.Is(err, comments.ErrNotAllowed) {
			return nil, ewrap.ErrNotAllowed
		}

		if errors.Is(err, comments.AccessDenied) {
			return nil, ewrap.AccessDenied
		}

		return nil, ewrap.InternalError
	}

	return &server.NewCommentResponse{
		Id:      commentID,
		Created: timestamppb.New(createdAt),
	}, nil
}

func ValidateContent(content string) error {
	if len(content) > 2000 {
		return ewrap.ErrMaxLength
	}
	return nil
}

// PostCommentToComment публикует комментарий к комментарию
func (s *ServerComments) PostCommentToComment(ctx context.Context, req *server.PostCommentRequest) (*server.NewCommentResponse, error) {
	const op = "internal/services/comments.PostCommentToComment"

	if err := ValidateContent(req.Content); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	commentID, createdAt, err := s.comments.PostCommentToComment(ctx, req.GetToken(), req.GetPostId(), req.GetParentId(), req.GetContent())
	if err != nil {
		if errors.Is(err, comments.ErrConnectionTime) {
			return nil, ewrap.ErrConnectionTime
		}

		if errors.Is(err, comments.AccessDenied) {
			return nil, ewrap.AccessDenied
		}

		return nil, ewrap.InternalError
	}

	return &server.NewCommentResponse{
		Id:      commentID,
		Created: timestamppb.New(createdAt),
	}, nil
}

// GetComments получает список комментариев к посту
func (s *ServerComments) GetComments(ctx context.Context, req *server.CommentsRequest) (*server.CommentsResponse, error) {
	const op = "internal/services/comments.GetComments"

	comment, err := s.comments.GetComments(ctx, req.GetPostId(), req.GetParentId())
	if err != nil {
		if errors.Is(err, comments.ErrGetComments) {
			return nil, ewrap.ErrGetComments
		}

		if errors.Is(err, comments.ErrFoundComments) {
			return nil, ewrap.ErrFoundComments
		}

		if errors.Is(err, comments.ErrConnectionTime) {
			return nil, ewrap.ErrConnectionTime
		}

		return nil, ewrap.InternalError
	}

	return ConvertToCommentsResponse(comment)
}

func ConvertToCommentsResponse(comments []models.Comment) (*server.CommentsResponse, error) {
	res := &server.CommentsResponse{}

	var serverComments []*server.Comment
	for _, comment := range comments {
		serverComments = append(serverComments, &server.Comment{
			Id:        comment.ID,
			UserId:    comment.UserID,
			PostId:    comment.PostID,
			Content:   comment.Content,
			CreatedAt: timestamppb.New(comment.CreatedAt),
			ParentId:  comment.ParentID,
		})
	}

	res.Comments = serverComments

	return res, nil
}
