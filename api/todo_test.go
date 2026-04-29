package api

import (
	"strings"
	"testing"
)

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
