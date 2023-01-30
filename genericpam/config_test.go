package genericpam

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestConfigureLogging(t *testing.T) {
	level := zerolog.GlobalLevel()
	defer zerolog.SetGlobalLevel(level)

	ConfigureLogging(LogConfig{Level: "trace"})

	assert.Equal(t, "trace", zerolog.GlobalLevel().String())
}

func TestConfigureLoggingError(t *testing.T) {
	level := zerolog.GlobalLevel()
	defer zerolog.SetGlobalLevel(level)

	ConfigureLogging(LogConfig{Level: "foo"})

	assert.Equal(t, "info", zerolog.GlobalLevel().String())
}
