package models

import "time"

// Post структура для хранения постов
type Post struct {
	/* идентификатор поста */
	ID int64 `json:"id"`
	/* идентификатор автора поста */
	UserID int64 `json:"user_id"`
	/* название поста */
	Title string `json:"title"`
	/* текст поста */
	Content string `json:"content"`
	/* время создания поста */
	CreatedAt time.Time `json:"created_at"`
	/* Указывает на то, что автор разрешил комментировать пост */
	AllowComments bool `json:"comments"`
}
