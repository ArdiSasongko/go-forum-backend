package middleware

import (
	"strconv"
	"time"

	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/env"
	"github.com/ArdiSasongko/go-forum-backend/internal/db/usersession"
	userrepository "github.com/ArdiSasongko/go-forum-backend/internal/repository/user.repository"
	"github.com/ArdiSasongko/go-forum-backend/pkg/database"
	"github.com/ArdiSasongko/go-forum-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func MiddlewareAuthValidate(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		logrus.WithField("get auth", "empty authorization header").Error("empty authorization")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "empty header authorization", nil)
	}

	dsn := env.GetEnv("DB_URL", "")
	db, err := database.InitDB(dsn)
	if err != nil {
		logrus.WithField("database", err.Error()).Fatal(err.Error())
	}

	token, err := userrepository.NewUserSessionRepository(db).GetTokenByToken(ctx.Context(), auth)
	if token == (usersession.UserSession{}) {
		logrus.WithField("get token", "token empty in database").Error("token empty in database")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "token empty in database", nil)
	} else if err != nil {
		logrus.WithField("get token", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "failed get token", nil)
	}

	claims, err := utils.ValidateToken(ctx.Context(), auth)
	if err != nil {
		logrus.WithField("validate token", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "failed validated token", nil)
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		logrus.WithField("validate token", "token has expired").Error("token has expired")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "token has expired", nil)
	}

	isValid := claims.IsValid
	ctx.Locals("username", claims.Username)
	ctx.Locals("email", claims.Email)
	ctx.Locals("role", claims.Role)
	ctx.Locals("is_valid", strconv.FormatBool(isValid))
	return ctx.Next()
}

func MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		logrus.WithField("get auth", "empty authorization header").Error("empty authorization")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "empty header authorization", nil)
	}

	claims, err := utils.ValidateToken(ctx.Context(), auth)
	if err != nil {
		logrus.WithField("validate token", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "failed validated token", nil)
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		logrus.WithField("validate token", "token has expired").Error("token has expired")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "token has expired", nil)
	}

	isValid := claims.IsValid
	ctx.Locals("username", claims.Username)
	ctx.Locals("email", claims.Email)
	ctx.Locals("role", claims.Role)
	ctx.Locals("is_valid", strconv.FormatBool(isValid))
	return ctx.Next()
}
