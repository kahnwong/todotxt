package api

import (
	"github.com/gofiber/fiber/v3"
)

type UpdateTodoRequest struct {
	ID      int    `json:"id" validate:"required"`
	Project string `json:"project"`
	Status  string `json:"status" validate:"required"`
	Context string `json:"context"`
}

type UpdateTodoTextRequest struct {
	ID   int    `json:"id" validate:"required"`
	Text string `json:"text" validate:"required"`
}

var todoService = &TodoService{}

func TodayController(c fiber.Ctx) error {
	return c.JSON(todoService.today())
}

func TinkeringController(c fiber.Ctx) error {
	return c.JSON(todoService.tinkering())
}

func WorkController(c fiber.Ctx) error {
	return c.JSON(todoService.work())
}

func UpdateTodoController(c fiber.Ctx) error {
	var req UpdateTodoRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := todoService.updateTodo(req.ID, req.Project, req.Status, req.Context)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Todo updated successfully"})
}

func UpdateTodoContentController(c fiber.Ctx) error {
	var req UpdateTodoTextRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := todoService.updateTodoContent(req.ID, req.Text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Todo text updated successfully"})
}
