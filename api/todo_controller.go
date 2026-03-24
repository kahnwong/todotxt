package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateTodoRequest struct {
	ID      int    `json:"id" binding:"required"`
	Project string `json:"project"`
	Status  string `json:"status" binding:"required"`
	Context string `json:"context"`
}

type UpdateTodoTextRequest struct {
	ID   int    `json:"id" binding:"required"`
	Text string `json:"text" binding:"required"`
}

var todoService = &TodoService{}

func TodayController(c *gin.Context) {
	c.JSON(http.StatusOK, todoService.today())
}

func TinkeringController(c *gin.Context) {
	c.JSON(http.StatusOK, todoService.tinkering())
}

func WorkController(c *gin.Context) {
	c.JSON(http.StatusOK, todoService.work())
}

func UpdateTodoController(c *gin.Context) {
	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := todoService.updateTodo(req.ID, req.Project, req.Status, req.Context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully"})
}

func UpdateTodoContentController(c *gin.Context) {
	var req UpdateTodoTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := todoService.updateTodoContent(req.ID, req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo text updated successfully"})
}
