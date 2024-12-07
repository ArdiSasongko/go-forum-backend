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
}

type ImageContent struct {
	ImageURL string `json:"image_url"`
}
