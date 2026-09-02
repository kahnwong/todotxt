package logging

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func configureLevel(t *testing.T, value *string) {
	t.Helper()

	original, existed := os.LookupEnv("LOG_LEVEL")
	t.Cleanup(func() {
		if existed {
			_ = os.Setenv("LOG_LEVEL", original)
		} else {
			_ = os.Unsetenv("LOG_LEVEL")
		}
		Configure()
	})

	if value == nil {
		_ = os.Unsetenv("LOG_LEVEL")
	} else {
		_ = os.Setenv("LOG_LEVEL", *value)
	}
	Configure()
}

func TestDefaultLogLevelIsDebug(t *testing.T) {
	configureLevel(t, nil)

	if level := zerolog.GlobalLevel(); level != zerolog.DebugLevel {
		t.Fatalf("expected debug level, got %s", level)
	}
	if !slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		t.Fatal("expected slog debug messages to be enabled")
	}
}

func TestLogLevelFromEnvironment(t *testing.T) {
	level := "error"
	configureLevel(t, &level)

	if got := zerolog.GlobalLevel(); got != zerolog.ErrorLevel {
		t.Fatalf("expected error level, got %s", got)
	}
	if slog.Default().Enabled(context.Background(), slog.LevelInfo) {
		t.Fatal("expected slog info messages to be disabled")
	}
}

func TestInvalidLogLevelFallsBackToDebug(t *testing.T) {
	level := "invalid"
	configureLevel(t, &level)

	if got := zerolog.GlobalLevel(); got != zerolog.DebugLevel {
		t.Fatalf("expected debug level, got %s", got)
	}
}
