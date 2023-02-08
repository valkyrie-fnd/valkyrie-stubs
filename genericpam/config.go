package genericpam

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Config struct {
	PamAPIKey      string            // pam API key
	ProviderTokens map[string]string // provider specific tokens used for reconciliation
	Address        string            // address to bind to, for example "localhost:8080"
	LogConfig      LogConfig         // logging configuration
}

type LogConfig struct {
	Level string // configured logging level
}

func ConfigureLogging(config LogConfig) {
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return pkgerrors.MarshalStack(errors.WithStack(err))
	}
	// configure stack field name to stack_trace, that is automatically recognized by Google error reporting
	zerolog.ErrorStackFieldName = "stack_trace"

	// use RFC3339 with nano precision for timestamp field
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// configure output to stdout
	log.Logger = log.Logger.Output(os.Stdout)

	// use global logger as default context logger (used when context is missing a logger: "zerolog.Ctx(ctx).Info()")
	zerolog.DefaultContextLogger = &log.Logger

	log.Info().Msg("Configured logging")
}
