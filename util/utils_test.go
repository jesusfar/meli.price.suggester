package util

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Log("NewLogger returns a Logger pointer")
	logger := NewLogger()
	assert.NotNil(t, logger)
	assert.Equal(t, LOG_INFO, logger.logLevel)
}

func TestLogger_Info(t *testing.T) {
	// Set LogLevel
	os.Setenv("LOG_LEVEL", "INFO")

	t.Log("Given LOG_LEVEL = INFO shows: ")
	{
		logger := NewLogger()
		assert.Equal(t, LOG_INFO, logger.logLevel)

		logger.Info("My log Info is showed.")
		logger.Debug("My log Debug is not showed.")
	}
}

func TestLogger_Debug(t *testing.T) {

	// Set LogLevel
	os.Setenv("LOG_LEVEL", "DEBUG")

	t.Log("Given LOG_LEVEL = DEBUG shows: ")
	{
		logger := NewLogger()
		assert.Equal(t, LOG_DEBUG, logger.logLevel)

		logger.Info("My log Info is showed.")
		logger.Warning("My log Warning is showed.")
		logger.Debug("My log Debug is showed.")
	}
}

func TestLogger_Warning(t *testing.T) {

	// Set LogLevel
	os.Setenv("LOG_LEVEL", "WARNING")

	t.Log("Given LOG_LEVEL = WARNING shows: ")
	{
		logger := NewLogger()
		assert.Equal(t, LOG_WARNING, logger.logLevel)

		logger.Info("My log Info is showed.")
		logger.Warning("My log Warning is showed.")
		logger.Debug("My log Debug is not showed.")
	}
}
