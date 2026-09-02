package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/kahnwong/todotxt/api"
	_ "github.com/kahnwong/todotxt/internal/logging"
	slogfiber "github.com/samber/slog-fiber"
	"github.com/sethvargo/go-envconfig"
)

var Config Env

type Env struct {
	ListenAddr string `env:"LISTEN_ADDR,default=:3000"`
}

//go:embed all:frontend/dist
var frontendAssets embed.FS

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func newApp(spaFS fs.FS) (*fiber.App, error) {
	index, err := fs.ReadFile(spaFS, "index.html")
	if err != nil {
		return nil, err
	}

	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})
	app.Use(slogfiber.New(slog.Default()))
	app.Use(recover.New())

	// API routes
	app.Get("/api/todo/today", api.TodayController)
	app.Get("/api/todo/tinkering", api.TinkeringController)
	app.Get("/api/todo/work", api.WorkController)
	app.Put("/api/todo/update", api.UpdateTodoController)
	app.Put("/api/todo/update-content", api.UpdateTodoContentController)

	app.Get("/*", static.New("", static.Config{
		FS: spaFS,
		Next: func(c fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api")
		},
		NotFoundHandler: func(c fiber.Ctx) error {
			c.Type("html")
			return c.Status(fiber.StatusOK).Send(index)
		},
	}))

	return app, nil
}

func main() {
	spaFS, err := fs.Sub(frontendAssets, "frontend/dist/spa")
	if err != nil {
		slog.Error("Failed to load frontend assets", "error", err)
		return
	}

	app, err := newApp(spaFS)
	if err != nil {
		slog.Error("Failed to initialize server", "error", err)
		return
	}

	if err := app.Listen(Config.ListenAddr); err != nil {
		slog.Error("Error starting server", "error", err)
	}
}

func init() {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &Config); err != nil {
		slog.Error("Failed to process environment variables", "error", err)
		os.Exit(1)
	}
}
