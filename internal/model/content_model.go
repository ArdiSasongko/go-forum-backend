package model

import (
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
)

type ContentModel struct {
	UserID         int32                   `json:"user_id"`
	Username       string                  `json:"username"`
	ContentTitle   string                  `json:"content_title" form:"content_title" validate:"required,min=1,max=255"`
	ContentBody    string                  `json:"content_body" form:"content_body" validate:"required"`
	ContentHastags []string                `json:"content_hastags" form:"content_hastags" validate:"omitempty"`
	Files          []*multipart.FileHeader `json:"files" validate:"omitempty"`
}

func (u ContentModel) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

// many
type ContentsResponse struct {
	ContentID      int            `json:"content_id"`
	ContentTitle   string         `json:"content_title"`
	ContentBody    string         `json:"content_body"`
	ContentImage   []ImageContent `json:"image_content"`
	ContentHastags []string       `json:"content_hastags"`
}

// one
type ContentResponse struct {
	ContentID      int            `json:"content_id"`
	ContentTitle   string         `json:"content_title"`
	ContentBody    string         `json:"content_body"`
	ContentImage   []ImageContent `json:"image_content"`
	ContentHastags []string       `json:"content_hastags"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	CreatedBy      string         `json:"created_by"`
	ContentMetrics `json:"content_metrics"`
}

type ContentMetrics struct {
	IsLike       bool `json:"is_liked"`
	LikeCount    int  `json:"like_count"`
	DislikeCount int  `json:"dislike_count"`
	CommentCount int  `json:"comment_count"`
	Pagination   `json:"pagination"`
	Comments     []CommentsResponse `json:"comments"`
}

type Pagination struct {
	Limit int32 `json:"limit"`
}
type CommentsResponse struct {
	Username  string    `json:"username"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ImageContent struct {
	ImageURL string `json:"image_url"`
}

type UpdateContent struct {
	ContentTitle   string   `json:"content_title" validate:"omitempty"`
	ContentBody    string   `json:"content_body" validate:"omitempty"`
	ContentHastags []string `json:"content_hastags" validate:"omitempty"`
	UpdatedBy      string
}

func (u UpdateContent) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

type CommentModel struct {
	UserID    int32  `json:"user_id"`
	ContentID int32  `json:"content_id"`
	Username  string `json:"username"`
	Comment   string `json:"comment" validate:"omitempty"`
}

func (u CommentModel) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

type LikedModel struct {
	UserID    int32  `json:"user_id"`
	ContentID int32  `json:"content_id"`
	Username  string `json:"username"`
	IsLike    bool   `json:"is_liked" validate:"omitempty"`
	IsDislike bool   `json:"is_disliked" validate:"omitempty"`
}

func (u LikedModel) Validate() error {
	v := validator.New()
	return v.Struct(u)
}
