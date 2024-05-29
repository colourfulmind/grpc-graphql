package tests

var resultSaveUserSuccess struct {
	Data struct {
		SaveUser struct {
			UserID int `json:"user_id"`
		} `json:"SaveUser"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultSaveUserFail struct {
	Data struct {
		SaveUser struct {
			UserID int `json:"user_id"`
		} `json:"SaveUser"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultLoginSuccess struct {
	Data struct {
		Login struct {
			Token string `json:"token"`
		} `json:"Login"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultLoginFail struct {
	Data struct {
		Login struct {
			Token string `json:"token"`
		} `json:"Login"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultSavePostSuccess struct {
	Data struct {
		SavePostResponse struct {
			PostID    int    `json:"post_id"`
			CreatedAt string `json:"created_at"`
		} `json:"SavePost"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultSavePostFail struct {
	Data struct {
		SavePostResponse struct {
			PostID    int    `json:"post_id"`
			CreatedAt string `json:"created_at"`
		} `json:"SavePost"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultProvidePostSuccess struct {
	Data struct {
		ProvidePostResponse struct {
			ID        int    `json:"id"`
			UserID    int    `json:"user_id"`
			Title     string `json:"title"`
			Content   string `json:"content"`
			CreatedAt string `json:"created_at"`
			Comments  bool   `json:"comments"`
		} `json:"ProvidePost"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultProvidePostFail struct {
	Data struct {
		ProvidePostResponse struct {
			ID        int    `json:"id"`
			UserID    int    `json:"user_id"`
			Title     string `json:"title"`
			Content   string `json:"content"`
			CreatedAt string `json:"created_at"`
			Comments  bool   `json:"comments"`
		} `json:"ProvidePost"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultProvideAllPosts struct {
	Data struct {
		ProvideAllPostsResponse struct {
			Posts []struct {
				ID        int    `json:"id"`
				UserID    int    `json:"user_id"`
				Title     string `json:"title"`
				Content   string `json:"content"`
				CreatedAt string `json:"created_at"`
				Comments  bool   `json:"comments"`
			} `json:"posts"`
		} `json:"ProvideAllPosts"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultSaveCommentSuccess struct {
	Data struct {
		SaveCommentResponse struct {
			ID        int    `json:"id"`
			CreatedAt string `json:"created_at"`
		} `json:"SaveComment"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultSaveCommentFail struct {
	Data struct {
		SaveCommentResponse struct {
			ID        int    `json:"id"`
			CreatedAt string `json:"created_at"`
		} `json:"SaveComment"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultSaveCommentToCommentSuccess struct {
	Data struct {
		SaveCommentResponse struct {
			ID        int    `json:"id"`
			CreatedAt string `json:"created_at"`
		} `json:"SaveCommentToComment"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultSaveCommentToCommentFail struct {
	Data struct {
		SaveCommentResponse struct {
			ID        int    `json:"id"`
			CreatedAt string `json:"created_at"`
		} `json:"SaveComment"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

var resultProvideComment struct {
	Data struct {
		ProvideCommentResponse struct {
			Comment []struct {
				ID        int    `json:"id"`
				UserID    int    `json:"user_id"`
				PostID    int    `json:"post_id"`
				Content   string `json:"content"`
				CreatedAt string `json:"created_at"`
				ParentID  int    `json:"parent_id"`
			} `json:"comments"`
		} `json:"ProvideComment"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}
