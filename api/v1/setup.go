package v1

import "github.com/gofiber/fiber/v2"

type Router interface {
	InstallRouter(app *fiber.App)
}

func InstallRouter(app *fiber.App) {
	setup(app, &ApiRouter{})
}

func setup(app *fiber.App, router ...Router) {
	for _, r := range router {
		r.InstallRouter(app)
	}
}
