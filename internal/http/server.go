package http

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
)

//go:embed views/*
var htmlFiles embed.FS

type Server struct {
	a *fiber.App
}

func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return s.a.Listen(fmt.Sprintf(":%s", port))
}

func (s *Server) Shutdown() error {
	return s.a.Shutdown()
}

func NewServer() *Server {
	viewEngine := html.NewFileSystem(http.FS(htmlFiles), ".html")

	a := fiber.New(fiber.Config{
		Views: viewEngine,
	})

	AddRoutes(a)

	return &Server{a}
}
