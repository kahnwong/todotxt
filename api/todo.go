package api

import (
	"fmt"
	"os"
	"strings"
	"time"

	todo "github.com/1set/todotxt"
)

type Todo struct {
	ID      int    `json:"id"`
	Context string `json:"context"`
	Project string `json:"project"`
	Todo    string `json:"todo"`
	Status  string `json:"status"`
}

func parseTodos(tasks todo.TaskList) []Todo {
	var todos []Todo

	for _, t := range tasks {
		// context
		var context string
		if len(t.Contexts) > 0 {
			context = fmt.Sprintf("@%s", t.Contexts[0])
		}

		// project
		var project string
		if len(t.Projects) > 0 {
			project = fmt.Sprintf("+%s", t.Projects[0])
		}

		// status - parse from todo text (e.g., =backlog, =in-progress, =done)
		status := "backlog" // default status
		words := strings.Fields(t.Todo)
		var cleanedWords []string

		for _, word := range words {
			if strings.HasPrefix(word, "=") {
				status = strings.TrimPrefix(word, "=")
			} else {
				cleanedWords = append(cleanedWords, word)
			}
		}

		// Rebuild todo text without status tag
		todoText := strings.Join(cleanedWords, " ")

		// append
		todos = append(todos, Todo{
			t.ID,
			context,
			project,
			todoText,
			status,
		})
	}

	return todos
}

type TodoService struct{}

func (ts *TodoService) loadTasklist() (todo.TaskList, error) {
	// Common setup - always executed
	sanitizedContent, err := sanitizeTodo() // strip leading `https://` which results in the todo body returning null
	if err != nil {
		fmt.Printf("Error reading todo.txt: %s", err)
		return nil, err
	}

	// Create in-memory pipe
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Printf("Error creating pipe: %s", err)
		return nil, err
	}

	// Write sanitized content to pipe in goroutine
	go func() {
		defer writer.Close()
		writer.WriteString(sanitizedContent)
	}()

	// Load from pipe
	tasklist, err := todo.LoadFromFile(reader)
	reader.Close()
	if err != nil {
		fmt.Printf("Error parsing todo.txt: %s", err)
		return nil, err
	}
	return tasklist, nil
}

func (ts *TodoService) filterTodos(filterFunc func(todo.TaskList) todo.TaskList) []Todo {
	var todos []Todo

	tasklist, err := ts.loadTasklist()
	if err != nil {
		return todos
	}

	filtered := filterFunc(tasklist)
	_ = filtered.Sort(todo.SortPriorityAsc, todo.SortProjectAsc)
	todos = parseTodos(filtered)

	return todos
}

func (ts *TodoService) today() []Todo {
	return ts.filterTodos(func(tl todo.TaskList) todo.TaskList {
		return tl.Filter(todo.FilterNotCompleted).Filter(todo.FilterDueToday, todo.FilterOverdue)
	})
}

func (ts *TodoService) tinkering() []Todo {
	return ts.filterTodos(func(tl todo.TaskList) todo.TaskList {
		return tl.Filter(todo.FilterNotCompleted).Filter(todo.FilterByContext("tinkering"))
	})
}

func (ts *TodoService) updateTodo(id int, newProject, newStatus, newContext string) error {
	// Read original file
	todoPath := os.Getenv("TODO_PATH")
	content, err := os.ReadFile(todoPath)
	if err != nil {
		return fmt.Errorf("error reading todo file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string
	currentID := 1

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			updatedLines = append(updatedLines, line)
			continue
		}

		// Check if this is the todo we want to update
		if currentID == id {
			// Parse and update the line
			updatedLine := ts.updateTodoLine(line, newProject, newStatus, newContext)
			updatedLines = append(updatedLines, updatedLine)
		} else {
			updatedLines = append(updatedLines, line)
		}
		currentID++
	}

	// Write back to file
	newContent := strings.Join(updatedLines, "\n")
	err = os.WriteFile(todoPath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing todo file: %w", err)
	}

	return nil
}

func (ts *TodoService) updateTodoContent(id int, newText string) error {
	// Read original file
	todoPath := os.Getenv("TODO_PATH")
	content, err := os.ReadFile(todoPath)
	if err != nil {
		return fmt.Errorf("error reading todo file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string
	currentID := 1

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			updatedLines = append(updatedLines, line)
			continue
		}

		// Check if this is the todo we want to update
		if currentID == id {
			updatedLines = append(updatedLines, newText)
		} else {
			updatedLines = append(updatedLines, line)
		}
		currentID++
	}

	// Write back to file
	newContent := strings.Join(updatedLines, "\n")
	err = os.WriteFile(todoPath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing todo file: %w", err)
	}

	return nil
}

func (ts *TodoService) updateTodoLine(line, newProject, newStatus, newContext string) string {
	words := strings.Fields(line)
	var updatedWords []string

	// Remove existing project, status, context, and due tags
	for _, word := range words {
		if !strings.HasPrefix(word, "+") && !strings.HasPrefix(word, "=") && !strings.HasPrefix(word, "@") && !strings.HasPrefix(word, "due:") {
			updatedWords = append(updatedWords, word)
		}
	}

	// Add new context tag if provided
	if newContext != "" {
		// Remove the @ prefix if it exists in newContext
		context := strings.TrimPrefix(newContext, "@")
		updatedWords = append(updatedWords, "@"+context)
	}

	// Add new project tag if provided
	if newProject != "" && newProject != "No Project" {
		// Remove the + prefix if it exists in newProject
		project := strings.TrimPrefix(newProject, "+")
		updatedWords = append(updatedWords, "+"+project)
	}

	// Add new status tag
	updatedWords = append(updatedWords, "="+newStatus)

	// Add due date for in-progress and done status
	if newStatus == "in-progress" || newStatus == "done" {
		today := time.Now().Format("2006-01-02")
		updatedWords = append(updatedWords, "due:"+today)
	}

	return strings.Join(updatedWords, " ")
}
