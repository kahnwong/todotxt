package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kahnwong/todotxt/api"

	"os"
)

func main() {
	// init
	router := gin.Default()

	// API routes
	router.GET("/api/todo/today", api.TodayController)
	router.GET("/api/todo/tinkering", api.TinkeringController)
	router.PUT("/api/todo/update", api.UpdateTodoController)
	router.PUT("/api/todo/update-content", api.UpdateTodoContentController)

	// Static routes
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// If the path doesn't start with /api, try to serve static files
		if !strings.HasPrefix(path, "/api") {
			c.File("/frontend/dist/spa" + path) // [TODO] dev mode use `frontend/dist/spa`
		}
	})

	// start server
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":3000"
	}
	err := router.Run(listenAddr)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
