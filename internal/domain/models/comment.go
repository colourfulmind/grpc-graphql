package models

import "time"

// Comment структура для хранения комментариев
type Comment struct {
	/* идентификатор комментария */
	ID int64 `json:"id"`
	/* идентификатор пользователя, отправившего комментарий */
	UserID int64 `json:"user_id"`
	/* идентификатор поста */
	PostID int64 `json:"post_id"`
	/* текст комментария */
	Content string `json:"content"`
	/* время создания комментария */
	CreatedAt time.Time `json:"created_at"`
	/* Родительский комментарий */
	ParentID int64 `json:"comment_id"`
}
