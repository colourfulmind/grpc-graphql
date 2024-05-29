package posts

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"ozon/internal/domain/models"
	"ozon/internal/ewrap"
	"ozon/internal/services/posts"
	server "ozon/protos/gen/go/posts"
)

const EmptyValue int64 = 0

// ServerPosts наследует интерфейс из protobuf и включает собственную реализацию
type ServerPosts struct {
	server.UnimplementedPostsServer
	posts Posts
}

// Posts - описыает методы для сервиса публикации постов
type Posts interface {
	PostNew(ctx context.Context, token string, title, content string, allowComments bool) (int64, time.Time, error)
	GetPostByID(ctx context.Context, postID int64) (models.Post, error)
	GetAllPosts(ctx context.Context, page int64) ([]models.Post, error)
}

// Register регистрирует сервер публикации постов
func Register(s *grpc.Server, posts Posts) {
	server.RegisterPostsServer(s, &ServerPosts{
		posts: posts,
	})
}

// PostNew публикует новый пост
// Возвращает ошибку, если у пользователя уже есть пост с таким же названием
// либо истекло время выполнения запроса
func (s *ServerPosts) PostNew(ctx context.Context, req *server.NewPostRequest) (*server.NewPostResponse, error) {
	const op = "internal/grpc/posts.PostNew"

	if err := ValidateArticle(req); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	postID, createdAt, err := s.posts.PostNew(ctx, req.GetToken(), req.GetTitle(), req.GetContent(), req.Comments)
	if err != nil {
		if errors.Is(err, posts.ErrPostExists) {
			return nil, fmt.Errorf("%s: %w", op, ewrap.PostAlreadyExists)
		}

		if errors.Is(err, posts.ErrConnectionTime) {
			return nil, fmt.Errorf("%s: %w", op, ewrap.ErrConnectionTime)
		}

		return nil, fmt.Errorf("%s: %w", op, ewrap.InternalError)
	}

	return &server.NewPostResponse{
		Id:        postID,
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}

// ValidateArticle проверяет, что заголовок и текст не пустые строки
func ValidateArticle(req *server.NewPostRequest) error {
	if req.GetTitle() == "" {
		return ewrap.TitleIsRequired
	}

	if req.GetContent() == "" {
		return ewrap.TextIsRequired
	}

	return nil
}

// GetPostByID получает пост по идентификатору
func (s *ServerPosts) GetPostByID(ctx context.Context, req *server.GetPostByIDRequest) (*server.GetPostByIDResponse, error) {
	const op = "internal/grpc/posts.GetPostByID"

	if err := ValidateID(req.Id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	post, err := s.posts.GetPostByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, posts.ErrPostNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ewrap.PostNotFound)
		}

		if errors.Is(err, posts.ErrConnectionTime) {
			return nil, fmt.Errorf("%s: %w", op, ewrap.ErrConnectionTime)
		}

		return nil, fmt.Errorf("%s: %w", op, ewrap.InternalError)
	}

	return &server.GetPostByIDResponse{
		Id:        post.ID,
		UserId:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: timestamppb.New(post.CreatedAt),
		Comments:  post.AllowComments,
	}, nil
}

func ValidateID(id int64) error {
	if id == EmptyValue {
		return ewrap.PostIDIsRequired
	}
	return nil
}

// GetAllPosts получает список всех постов
func (s *ServerPosts) GetAllPosts(ctx context.Context, req *server.GetAllPostsRequest) (*server.GetPostResponse, error) {
	const op = "internal/grpc/posts.GetAllPosts"

	post, err := s.posts.GetAllPosts(ctx, req.GetPage())
	if err != nil {
		if errors.Is(err, posts.ErrConnectionTime) {
			return nil, fmt.Errorf("%s: %w", op, ewrap.ErrConnectionTime)
		}

		return nil, fmt.Errorf("%s: %w", op, ewrap.InternalError)
	}

	return ConvertToPostResponse(post)
}

func ConvertToPostResponse(posts []models.Post) (*server.GetPostResponse, error) {
	res := &server.GetPostResponse{}

	var serverPosts []*server.Post
	for _, post := range posts {
		serverPosts = append(serverPosts, &server.Post{
			Id:        post.ID,
			UserId:    post.UserID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: timestamppb.New(post.CreatedAt),
			Comments:  post.AllowComments,
		})
	}

	res.Posts = serverPosts

	return res, nil
}
