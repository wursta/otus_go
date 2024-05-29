package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigSuccess(t *testing.T) {
	config, err := NewConfig("./test/config_test_ok.toml")

	require.Nil(t, err)
	require.Equal(t, "INFO", config.Logger.Level)
	require.Equal(t, "localhost", config.Server.Host)
	require.Equal(t, "8080", config.Server.Port)
}

func TestNewConfigError(t *testing.T) {
	t.Run("unknown file", func(t *testing.T) {
		_, err := NewConfig("./test/unknown_config.toml")
		require.NotNil(t, err)
	})

	t.Run("file invalid format", func(t *testing.T) {
		_, err := NewConfig("./test/config_test_fail.toml")
		require.NotNil(t, err)
	})
}
