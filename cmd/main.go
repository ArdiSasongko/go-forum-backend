package main

import (
	"fmt"

	v1 "github.com/ArdiSasongko/go-forum-backend/api/v1"
	"github.com/ArdiSasongko/go-forum-backend/env"
	contentservice "github.com/ArdiSasongko/go-forum-backend/internal/service/content.service"
	userservice "github.com/ArdiSasongko/go-forum-backend/internal/service/user.service"
	"github.com/ArdiSasongko/go-forum-backend/pkg/database"
	"github.com/ArdiSasongko/go-forum-backend/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Setup(logger *logrus.Logger) *fiber.App {
	env.SetupEnvFile()
	dsn := env.GetEnv("DB_URL", "")
	db, err := database.InitDB(dsn)
	if err != nil {
		logger.WithError(err)
	}

	contentService := contentservice.NewContentService(db, logger)
	userService := userservice.NewUserService(db, logger)
	apiRouter := v1.NewApiRouter(userService, contentService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	// app.Use(func(ctx *fiber.Ctx) error {
	// 	return ctx.SendStatus(fiber.StatusNotFound)
	// })

	apiRouter.InstallRouter(app)
	return app
}

func main() {
	logger := log.InitLogger()
	app := Setup(logger)
	PORT := env.GetEnv("APP_PORT", "4000")
	HOST := env.GetEnv("APP_HOST", "0.0.0.0")
	logger.WithField("port", PORT).Info("Starting Server")
	logger.Fatal(app.Listen(fmt.Sprintf("%s:%s", HOST, PORT)))
}
