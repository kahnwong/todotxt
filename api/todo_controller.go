package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var todoService = &TodoService{}

func TodayController(c *gin.Context) {
	c.JSON(http.StatusOK, todoService.today())
}

func TinkeringController(c *gin.Context) {
	c.JSON(http.StatusOK, todoService.tinkering())
}
