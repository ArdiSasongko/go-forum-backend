package main

import (
	"fmt"

	"github.com/ArdiSasongko/go-forum-backend/api/router"
	"github.com/ArdiSasongko/go-forum-backend/env"
	"github.com/ArdiSasongko/go-forum-backend/internal/repository"
	"github.com/ArdiSasongko/go-forum-backend/internal/service"
	"github.com/ArdiSasongko/go-forum-backend/pkg/database"
	"github.com/gofiber/fiber/v2"
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

	userRepo := repository.NewuserRepository(db)
	userSessionRepo := repository.NewUserSessionRepository(db)
	userService := service.NewUserService(userRepo, userSessionRepo, db)
	apiRouter := router.NewApiRouter(userService)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	apiRouter.InstallRouter(app)
	return app
}

func main() {
	app := Setup()
	logrus.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_PORT", "4000"))))
}
