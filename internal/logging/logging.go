package logging

import (
	"log/slog"
	"os"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

const defaultLevel = zerolog.InfoLevel

func init() {
	Configure()
}

func Configure() {
	level := defaultLevel
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		parsedLevel, err := zerolog.ParseLevel(envLevel)
		if err == nil {
			level = parsedLevel
		}
	}

	zerolog.SetGlobalLevel(level)
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	logger := zerolog.New(output).Level(level).With().Timestamp().Logger()
	slog.SetDefault(slog.New(slogzerolog.Option{Logger: &logger}.NewZerologHandler()))
}
