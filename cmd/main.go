package main

import (
	"fmt"

	v1 "github.com/ArdiSasongko/go-forum-backend/api/v1"
	"github.com/ArdiSasongko/go-forum-backend/env"
	contentservice "github.com/ArdiSasongko/go-forum-backend/internal/service/content.service"
	userservice "github.com/ArdiSasongko/go-forum-backend/internal/service/user.service"
	"github.com/ArdiSasongko/go-forum-backend/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Setup() *fiber.App {
	env.SetupEnvFile()
	dsn := env.GetEnv("DB_URL", "")
	db, err := database.InitDB(dsn)
	if err != nil {
		logrus.WithField("database", err.Error()).Fatal(err.Error())
	}

	contentService := contentservice.NewContentService(db)
	userService := userservice.NewUserService(db)
	apiRouter := v1.NewApiRouter(userService, contentService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New())

	apiRouter.InstallRouter(app)
	return app
}

func main() {
	app := Setup()
	logrus.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "0.0.0.0"), env.GetEnv("APP_PORT", "4000"))))
}
