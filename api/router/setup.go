package router

import "github.com/gofiber/fiber/v2"

type Router interface {
	InstallRouter(app *fiber.App)
}

func InstallRouter(app *fiber.App) {
	setUp(app, ApiRouter{})
}

func setUp(app *fiber.App, router ...Router) {
	for _, r := range router {
		r.InstallRouter(app)
	}
}
