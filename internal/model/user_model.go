package model

import "github.com/go-playground/validator/v10"

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
