# Agents Guide for todotxt Project

This document provides guidance for AI agents and developers using AI tools when working with this codebase.

## Project Overview

A web-based UI and API for managing and viewing todo.txt files. The application provides two main interfaces:

1. **List View** - Displays filtered tasks in two sections: "Today" (tasks due today or overdue) and "Tinkering" (tasks with @tinkering context)
2. **Kanban Board** - Drag-and-drop board with swimlanes organized by project, featuring three status columns (Backlog, In Progress, Done) with inline editing capabilities

## Technology Stack

### Backend
- **Language**: Go 1.25.4
- **Framework**: Gin web framework (v1.11.0)
- **Key Libraries**:
  - `github.com/1set/todotxt` (v0.0.4) - todo.txt parsing
- **Environment**: direnv with .envrc (loads .env via dotenv)

### Frontend
- **Framework**: Vue 3 (v3.4.18) with TypeScript
- **UI Library**: Quasar Framework (v2.16.0) - Material Design components
- **Build Tool**: Vite
- **Package Manager**: Yarn
- **HTTP Client**: Axios
- **Routing**: Vue Router
- **Icons**: @iconify/vue (v5.0.0)

### Infrastructure
- Docker with multi-stage builds (golang:1.25-alpine → distroless)
- GitHub Actions for CI/CD
- Pre-commit hooks for code quality
- SOPS for secret encryption (.sops.yaml with age key)

## Project Structure

```
/
├── api/                    # Go backend API
│   ├── todo.go            # Core todo parsing and filtering logic
│   ├── todo_controller.go # HTTP request handlers
│   └── utils.go           # Utility functions (URL sanitization)
├── frontend/              # Vue.js frontend application
│   ├── src/
│   │   ├── components/    # Vue components (TodoComponent.vue)
│   │   ├── pages/         # Page components (IndexPage.vue, KanbanPage.vue)
│   │   ├── layouts/       # Layout components (MainLayout.vue)
│   │   └── router/        # Vue Router configuration
│   └── dist/spa/          # Built frontend (embedded in Go binary)
├── main.go                # Application entry point
├── Dockerfile             # Multi-stage Docker build
├── Makefile               # Build and development commands
└── .github/workflows/     # CI/CD pipeline definitions
```

## Key Files and Their Purposes

### Backend (Go)
- `main.go` - Application entry point, sets up Gin router and embeds frontend
- `api/todo.go` - Core business logic for parsing and filtering todo.txt files
- `api/todo_controller.go` - HTTP controllers for `/api/todo/today` and `/api/todo/tinkering`
- `api/utils.go` - URL sanitization to prevent parsing issues

### Frontend (Vue)
- `frontend/src/pages/IndexPage.vue` - List view with tab navigation for Today and Tinkering sections
- `frontend/src/pages/KanbanPage.vue` - Kanban board with drag-and-drop, swimlanes by project, inline editing
- `frontend/src/components/TodoComponent.vue` - Reusable component for displaying todo lists
- `frontend/src/router/routes.ts` - Route definitions
- `frontend/src/layouts/MainLayout.vue` - Application layout wrapper

## API Endpoints

- `GET /api/todo/today` - Returns tasks due today or overdue
- `GET /api/todo/tinkering` - Returns tasks with @tinkering context
- `PUT /api/todo/update` - Update task project, status, and context
- `PUT /api/todo/update-text` - Update full task text

Response format:
```json
[
  {
    "id": 1,
    "context": "@work",
    "project": "+project-name",
    "todo": "Task description",
    "created_date": "2026-03-15"
  }
]
```

## Frontend Routes

- `/` - IndexPage (list view with Today/Tinkering sections)
- `/kanban` - KanbanPage (kanban board view)
- `/*` - ErrorNotFound (404 page)

## Configuration

Environment variables (`.env` file):
- `LISTEN_ADDR` - Server listen address (default: `:3000`)
- `TODO_PATH` - Path to todo.txt file (required)

## Development Workflow

### Starting Development Servers
```bash
make start-frontend  # Start Vue dev server (port 9000)
make start-backend   # Run Go backend (port 3000)
```

### Building
```bash
make build-frontend  # Build Vue app for production
make build          # Build Go binary with embedded frontend
```

### Testing
```bash
make test           # Run Go tests
```

### Code Quality
- Pre-commit hooks automatically run:
  - `go fmt` - Code formatting
  - `golangci-lint` - Linting
  - `go test` - Unit tests

## Common Tasks

### Adding a New Todo Filter View
1. Add filter function in `api/todo.go` (follow pattern of `GetTodayTodos` or `GetTinkeringTodos`)
2. Add controller endpoint in `api/todo_controller.go`
3. Register route in `main.go`
4. Add frontend API call in `frontend/src/pages/IndexPage.vue`
5. Add new tab in the Quasar tabs component

### Modifying Todo Display
- Edit `frontend/src/components/TodoComponent.vue`
- Context tags are styled with blue color
- Project tags are styled with gold color

### Changing URL Sanitization
- Modify `api/utils.go` - currently strips `https://` prefix to prevent parsing issues

## Code Conventions

### Go
- Follow standard Go formatting (enforced by `go fmt`)
- Use descriptive variable names
- Keep functions focused and single-purpose
- Error handling: return errors, don't panic

### Vue/TypeScript
- Use Composition API with `<script setup lang="ts">`
- TypeScript for type safety
- Quasar components for UI consistency
- Props and emits should be typed

## Deployment

### Docker
```bash
docker build -t todotxt .
docker run -p 3000:3000 -v /path/to/todo.txt:/todo.txt -e TODO_PATH=/todo.txt todotxt
```

### CI/CD
- GitHub Actions automatically:
  - Run tests on PRs
  - Build and push Docker images on master branch
  - Images pushed to GitHub Container Registry

## Important Implementation Details

1. **Custom Status System**: The application extends todo.txt format with custom status tags:
   - `=backlog` - Default status (not written to file if default)
   - `=in-progress` - Task currently being worked on
   - `=done` - Completed task
   - Status tags are parsed from todo text and stripped from display
   - Due dates are automatically added when status changes to in-progress or done

2. **Frontend Embedding**: The built Vue app is embedded in the Go binary using `//go:embed` directive

3. **URL Sanitization**: URLs in todo.txt are sanitized (https:// stripped) to prevent parsing issues with the todotxt library

4. **Auto-refresh**: Frontend polls API every 10 seconds for updates

5. **Task Sorting**: Tasks are sorted by priority and project name

6. **Distroless Container**: Production image uses distroless base for minimal attack surface

7. **Kanban Board Logic**:
   - "Today" swimlane excludes tasks with both @tinkering AND a project
   - Project swimlanes created dynamically from @tinkering tasks
   - Drag-and-drop automatically updates status, project, and context
   - Moving from Today to project lane adds @tinkering context

## Testing Strategy

- **API Testing**: Hurl for HTTP API testing (hurl/today.hurl)
- **Backend**: Go unit tests for backend logic
- **Frontend**: No unit tests (test script exits successfully)
- **Pre-commit**: Hooks ensure tests pass before commit
- **CI/CD**: Pipeline runs tests on all PRs

## Troubleshooting

### Frontend not loading
- Ensure `make build-frontend` was run before building Go binary
- Check that `frontend/dist/spa` directory exists

### Todo.txt not parsing correctly
- Check URL sanitization in `api/utils.go`
- Verify todo.txt format follows standard conventions
- Check `TODO_PATH` environment variable is set correctly

### Port conflicts
- Frontend dev server uses port 9000
- Backend uses port 3000 (configurable via `LISTEN_ADDR`)

## Future Enhancement Ideas

- Add more filter views (by priority, by project, by context)
- Implement task editing/creation via API
- Add authentication/authorization
- Support multiple todo.txt files
- Add task completion tracking
- Implement search functionality
