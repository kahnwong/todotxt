package api

import (
	"strings"
	"testing"

	todo "github.com/1set/todotxt"
)

func TestParseTodosIncludesPriority(t *testing.T) {
	task, err := todo.ParseTask("(A) 2026-04-29 finish report @work +ops due:2026-04-29")
	if err != nil {
		t.Fatalf("failed to parse task: %v", err)
	}

	parsed := parseTodos(todo.TaskList{*task})
	if len(parsed) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(parsed))
	}

	if parsed[0].Priority != "A" {
		t.Fatalf("expected priority A, got %q", parsed[0].Priority)
	}
}

func TestUpdateTodoLinePreservesPriority(t *testing.T) {
	ts := &TodoService{}

	updated := ts.updateTodoLine("(A) 2026-04-29 finish report due:2026-04-29", "", "stuck", "")

	if !strings.HasPrefix(updated, "(A) 2026-04-29") {
		t.Fatalf("expected priority to be preserved near the start, got %q", updated)
	}
}

func TestUpdateTodoLinePreservesDueDateForStuck(t *testing.T) {
	ts := &TodoService{}

	updated := ts.updateTodoLine("2026-04-29 finish report due:2026-04-29", "", "stuck", "")

	if !strings.Contains(updated, "=stuck") {
		t.Fatalf("expected stuck status, got %q", updated)
	}

	if !strings.Contains(updated, "due:2026-04-29") {
		t.Fatalf("expected due date to be preserved, got %q", updated)
	}
}

func TestUpdateTodoLineRefreshesDueDateForInProgress(t *testing.T) {
	ts := &TodoService{}

	updated := ts.updateTodoLine("2026-04-01 finish report due:2026-04-15", "", "in-progress", "")

	if !strings.Contains(updated, "=in-progress") {
		t.Fatalf("expected in-progress status, got %q", updated)
	}

	if !strings.Contains(updated, "due:") {
		t.Fatalf("expected due date, got %q", updated)
	}

	if strings.Contains(updated, "due:2026-04-15") {
		t.Fatalf("expected due date to be refreshed, got %q", updated)
	}
}
