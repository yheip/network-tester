package http

import (
	"embed"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"github.com/gofiber/fiber/v3/middleware/filesystem"
	"github.com/yheip/network-tester/internal/curl"
)

//go:embed public/*
var publicDir embed.FS

func AddRoutes(a *fiber.App) {
	a.Use("/public", filesystem.New(filesystem.Config{
		Root:       publicDir,
		PathPrefix: "public",
	}))

	a.Use(encryptcookie.New(encryptcookie.Config{
		Key: "486J6n0zGiw6ygJ0562Z429LMlyv16f8",
	}))

	a.Get("/", func(c fiber.Ctx) error {
		return c.Render("views/index", fiber.Map{})
	})

	a.Post("/httptest", HandleHttpTest(curl.Run))
	a.Get("/sse", HandleSSE)
}
