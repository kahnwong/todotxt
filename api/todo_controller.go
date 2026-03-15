package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TodayController(c *gin.Context) {
	c.JSON(http.StatusOK, getTodos("today"))
}

func TinkeringController(c *gin.Context) {
	c.JSON(http.StatusOK, getTodos("tinkering"))
}
