package api

import (
	"fmt"

	todo "github.com/1set/todotxt"
)

var todotxtSanitizedPath = "/tmp/todo.txt"

type Todo struct {
	ID      int    `json:"id"`
	Context string `json:"context"`
	Project string `json:"project"`
	Todo    string `json:"todo"`
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

		// append
		todos = append(todos, Todo{
			t.ID,
			context,
			project,
			t.Todo,
		})
	}

	return todos
}

type TodoService struct{}

func (ts *TodoService) loadTasklist() (todo.TaskList, error) {
	// Common setup - always executed
	sanitizeTodo() // strip leading `https://` which results in the todo body returning null
	tasklist, err := todo.LoadFromPath(todotxtSanitizedPath)
	if err != nil {
		fmt.Printf("Error reading todo.txt: %s", err)
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
