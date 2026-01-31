package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func TodoTodayController(c *gin.Context) {
	c.JSON(http.StatusOK, getTodos("today"))
}

func TodoTinkeringController(c *gin.Context) {
	c.JSON(http.StatusOK, getTodos("tinkering"))
}
