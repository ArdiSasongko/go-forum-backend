package middleware

import (
	"strconv"
	"time"

	"github.com/ArdiSasongko/go-forum-backend/api/types"
	"github.com/ArdiSasongko/go-forum-backend/env"
	"github.com/ArdiSasongko/go-forum-backend/internal/sqlc/usersession"
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
	tx, err := db.BeginTx(ctx.Context(), nil)
	defer utils.Tx(tx, err)

	userSessionQueries := usersession.New(db).WithTx(tx)
	token, err := userSessionQueries.GetTokenByToken(ctx.Context(), auth)
	if token == (usersession.UserSession{}) {
		logrus.WithField("get token", "token is invalid").Error("token is invalid")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "token is invalid, not found in database", nil)
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
	ctx.Locals("user_id", claims.UserID)
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

	claims, err := utils.ValidateRefreshToken(ctx.Context(), auth)
	if err != nil {
		logrus.WithField("validate token", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "failed validated token", nil)
	}

	isValid := claims.IsValid
	ctx.Locals("user_id", claims.UserID)
	ctx.Locals("username", claims.Username)
	ctx.Locals("email", claims.Email)
	ctx.Locals("role", claims.Role)
	ctx.Locals("is_valid", strconv.FormatBool(isValid))
	return ctx.Next()
}
