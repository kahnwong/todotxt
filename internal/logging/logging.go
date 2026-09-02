package logging

import (
	"log/slog"
	"os"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

const defaultLevel = zerolog.DebugLevel

func init() {
	Configure()
}

func Configure() {
	level := defaultLevel
	var parseErr error
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		parsedLevel, err := zerolog.ParseLevel(envLevel)
		if err == nil {
			level = parsedLevel
		} else {
			parseErr = err
		}
	}

	zerolog.SetGlobalLevel(level)
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	logger := zerolog.New(output).Level(level).With().Timestamp().Logger()
	slog.SetDefault(slog.New(slogzerolog.Option{
		Level:  slogzerolog.ZeroLogLeveler{Logger: &logger},
		Logger: &logger,
	}.NewZerologHandler()))
	if parseErr != nil {
		slog.Warn("Invalid LOG_LEVEL; using debug", "value", os.Getenv("LOG_LEVEL"), "error", parseErr)
	}
}
