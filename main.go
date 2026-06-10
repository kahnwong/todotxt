package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kahnwong/todotxt/api"
	_ "github.com/kahnwong/todotxt/internal/logging"
	"github.com/sethvargo/go-envconfig"
)

var Config Env

type Env struct {
	ListenAddr string `env:"LISTEN_ADDR,default=:3000"`
	Mode       string `env:"MODE,default=development"`
}

func main() {
	// init
	router := gin.Default()

	// API routes
	router.GET("/api/todo/today", api.TodayController)
	router.GET("/api/todo/tinkering", api.TinkeringController)
	router.GET("/api/todo/work", api.WorkController)
	router.PUT("/api/todo/update", api.UpdateTodoController)
	router.PUT("/api/todo/update-content", api.UpdateTodoContentController)

	// Static routes
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// If the path doesn't start with /api, try to serve static files
		if !strings.HasPrefix(path, "/api") {
			staticPath := "/frontend/dist/spa" // for docker
			if Config.Mode == "development" {
				staticPath = "frontend/dist/spa"
			}
			c.File(staticPath + path)
		}
	})

	// start server
	err := router.Run(Config.ListenAddr)
	if err != nil {
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
