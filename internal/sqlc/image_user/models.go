// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package imageuser

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type Roles string

const (
	RolesAdmin Roles = "admin"
	RolesUser  Roles = "user"
)

func (e *Roles) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Roles(s)
	case string:
		*e = Roles(s)
	default:
		return fmt.Errorf("unsupported scan type for Roles: %T", src)
	}
	return nil
}

type NullRoles struct {
	Roles Roles
	Valid bool // Valid is true if Roles is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRoles) Scan(value interface{}) error {
	if value == nil {
		ns.Roles, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Roles.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRoles) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Roles), nil
}

type TokenType string

const (
	TokenTypeEmail         TokenType = "email"
	TokenTypePasswordReset TokenType = "password_reset"
)

func (e *TokenType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TokenType(s)
	case string:
		*e = TokenType(s)
	default:
		return fmt.Errorf("unsupported scan type for TokenType: %T", src)
	}
	return nil
}

type NullTokenType struct {
	TokenType TokenType
	Valid     bool // Valid is true if TokenType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTokenType) Scan(value interface{}) error {
	if value == nil {
		ns.TokenType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TokenType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTokenType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TokenType), nil
}

type Content struct {
	ID             int32
	UserID         int32
	ContentTitle   string
	ContentBody    string
	ContentHastags string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      string
	UpdatedBy      string
}

type ImagesContent struct {
	ID        int32
	ContentID int32
	ImageUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ImagesUser struct {
	ID        int32
	UserID    int32
	ImageUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Token struct {
	UserID    int32
	TokenType TokenType
	Token     int32
	ExpiredAt time.Time
}

type User struct {
	ID        int32
	Name      string
	Username  string
	Email     string
	Password  string
	Role      Roles
	IsValid   sql.NullBool
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

type UserSession struct {
	UserID              int32
	Token               string
	TokenExpired        time.Time
	RefreshToken        string
	RefreshTokenExpired time.Time
}
