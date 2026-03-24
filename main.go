package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kahnwong/todotxt/api"
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
		fmt.Println("Error starting server", err)
	}
}

func init() {
	ctx := context.Background()

	if err := envconfig.Process(ctx, &Config); err != nil {
		log.Fatal(err)
	}
}
