package middleware

import (
	"strconv"
	"strings"
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
		logrus.WithField("get auth", "empty authorization header").Error("UNAUTHORIZED")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "empty header authorization")
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
		logrus.WithField("get token", "token is invalid").Error("UNAUTHORIZED")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "token is invalid")
	} else if err != nil {
		logrus.WithField("get token", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	}

	claims, err := utils.ValidateToken(ctx.Context(), auth)
	if err != nil {
		logrus.WithField("validate token", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		logrus.WithField("validate token", "token has expired").Error("token has expired")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "token has expired")
	}

	isValid := claims.IsValid
	ctx.Locals("user_id", claims.UserID)
	ctx.Locals("username", claims.Username)
	ctx.Locals("email", claims.Email)
	ctx.Locals("role", claims.Role)
	ctx.Locals("is_valid", strconv.FormatBool(isValid))
	return ctx.Next()
}

func CheckValidUser(ctx *fiber.Ctx) error {
	IsValid := ctx.Locals("is_valid").(string)
	IsValid = strings.ToLower(IsValid)

	if IsValid != "true" {
		logrus.WithField("check valid", "please validation email").Error("unauthorized")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "Please Validation Your Email")
	}

	return ctx.Next()
}

func MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		logrus.WithField("get auth", "empty authorization header").Error("UNAUTHORIZED")
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "empty header authorization")
	}

	claims, err := utils.ValidateRefreshToken(ctx.Context(), auth)
	if err != nil {
		logrus.WithField("validate token", err.Error()).Error(err.Error())
		return types.SendResponse(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "failed validated token")
	}

	isValid := claims.IsValid
	ctx.Locals("user_id", claims.UserID)
	ctx.Locals("username", claims.Username)
	ctx.Locals("email", claims.Email)
	ctx.Locals("role", claims.Role)
	ctx.Locals("is_valid", strconv.FormatBool(isValid))
	return ctx.Next()
}
