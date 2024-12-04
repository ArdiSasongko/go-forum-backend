package model

import (
	"mime/multipart"

	"github.com/go-playground/validator/v10"
)

type UserModel struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Username string `json:"username" validate:"required,min=6,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
	Role     string `json:"role"`
}

func (u UserModel) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

func (u LoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

type ResponseLogin struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type PayloadToken struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type RefreshToken struct {
	Token string `json:"token" validate:"required"`
}

func (u RefreshToken) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

type ValidateToken struct {
	Token int32 `json:"token" validate:"required"`
}

func (u ValidateToken) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

type ValidatePayload struct {
	Token    int32  `json:"token"`
	Username string `json:"username"`
}

type SendEmail struct {
	Email string `json:"email" validate:"required,email"`
}

func (u SendEmail) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

type ResetPassword struct {
	Token           int32  `json:"token" validate:"required,number,min=6"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=255"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=255,eqfield=Password"`
}

func (u ResetPassword) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

type ProfileModel struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ImageURL string `json:"image_url"`
	IsValid  bool   `json:"is_valid"`
	Role     string `json:"role"`
}

type UpdateProfile struct {
	Email string                  `json:"email"`
	Files []*multipart.FileHeader `json:"file" validate:"omitempty"`
}

func (u UpdateProfile) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

type UpdateUser struct {
	Username string `json:"username" validate:"required,min=6,max=255"`
	Name     string `json:"name" validate:"omitempty,min=1,max=255"`
}

func (u UpdateUser) Validate() error {
	v := validator.New()
	return v.Struct(u)
}
