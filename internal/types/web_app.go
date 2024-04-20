package types

import "time"

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Dob       string `json:"dob"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
}

type EditUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Dob       *string `json:"dob"`
	Password  *string `json:"password"`
}

type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type CreatePostRequest struct {
	UserId           int64  `json:"user_id"`
	ContentText      string `json:"content_text"`
	ContentImagePath string `json:"content_image_path"`
	Visible          bool   `json:"visible"`
}

type EditPostRequest struct {
	ContentText      *string `json:"content_text"`
	ContentImagePath *string `json:"content_image_path"`
	Visible          *bool   `json:"visible"`
}

type CreatePostCommentRequest struct {
	ContentText string `json:"content_text"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type PostDetailResponse struct {
	PostID           int64     `json:"post_id"`
	UserID           int64     `json:"user_id"`
	ContentText      string    `json:"content_text"`
	ContentImagePath string    `json:"content_image_path"`
	Visible          bool      `json:"visible"`
	CreatedTime      time.Time `json:"created_time"`
}
