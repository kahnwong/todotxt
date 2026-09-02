package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gofiber/fiber/v3"
)

func testApp(t *testing.T) *fiber.App {
	t.Helper()

	app, err := newApp(fstest.MapFS{
		"index.html":    &fstest.MapFile{Data: []byte("<html>todo.txt</html>")},
		"assets/app.js": &fstest.MapFile{Data: []byte("console.log('todo.txt')")},
	})
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}
	return app
}

func TestStaticAsset(t *testing.T) {
	app := testApp(t)
	response, err := app.Test(httptest.NewRequest("GET", "/assets/app.js", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}
	if string(body) != "console.log('todo.txt')" {
		t.Fatalf("unexpected body %q", body)
	}
}

func TestSPAFallback(t *testing.T) {
	app := testApp(t)
	response, err := app.Test(httptest.NewRequest("GET", "/kanban", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}
	if string(body) != "<html>todo.txt</html>" {
		t.Fatalf("unexpected body %q", body)
	}
}

func TestUnknownAPIRouteDoesNotServeSPA(t *testing.T) {
	app := testApp(t)
	response, err := app.Test(httptest.NewRequest("GET", "/api/unknown", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != 404 {
		t.Fatalf("expected status 404, got %d", response.StatusCode)
	}
}

func TestUpdateTodoValidatesRequiredFields(t *testing.T) {
	app := testApp(t)
	request := httptest.NewRequest("PUT", "/api/todo/update", strings.NewReader(`{"id":1}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != 400 {
		t.Fatalf("expected status 400, got %d", response.StatusCode)
	}
}
