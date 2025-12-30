package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/kahnwong/todotxt/api"

	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

//go:embed frontend/dist/spa
var embedFS embed.FS // 2. Embed the 'spa' directory

func main() {
	// init
	app := fiber.New()

	// render site
	app.Get("/api/todo/today", api.TodoTodayController)
	app.Get("/api/todo/tinkering", api.TodoTinkeringController)
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(embedFS),
		PathPrefix: "frontend/dist/spa",
		Index:      "index.html",
	}))

	// start server
	err := app.Listen(os.Getenv("LISTEN_ADDR"))
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
