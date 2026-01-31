package main

import (
	"embed"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kahnwong/todotxt/api"

	"io/fs"
	"net/http"
	"os"
	"strings"
)

//go:embed frontend/dist/spa
var embedFS embed.FS

func main() {
	// init
	router := gin.Default()

	// API routes
	router.GET("/api/todo/today", api.TodoTodayController)
	router.GET("/api/todo/tinkering", api.TodoTinkeringController)

	// Static routes
	subFS, _ := fs.Sub(embedFS, "frontend/dist/spa")
	staticServer := http.FS(subFS)
	fileServer := http.FileServer(staticServer)

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		filePath := strings.TrimPrefix(path, "/")

		if filePath == "" {
			filePath = "index.html"
		}

		file, err := subFS.Open(filePath)
		if err == nil {
			file.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		c.FileFromFS("index.html", staticServer)
	})

	// start server
	err := router.Run(os.Getenv("LISTEN_ADDR"))
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
