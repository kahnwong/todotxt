package api

import (
	"os"
	"strings"
)

func sanitizeTodo() (string, error) {
	// Read and sanitize in-memory
	todoBytes, err := os.ReadFile(Config.TodoPath)
	if err != nil {
		return "", err
	}

	todoStr := string(todoBytes)
	todoSanitized := strings.ReplaceAll(todoStr, "https://", "")

	return todoSanitized, nil
}
