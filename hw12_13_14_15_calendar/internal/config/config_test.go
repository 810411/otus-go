package config //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Config(t *testing.T) {
	t.Run("read config", func(t *testing.T) {
		got, err := ParseConfig("../../configs/local.json")

		require.NoError(t, err)
		require.NotEmpty(t, got.HTTP)
		require.NotEmpty(t, got.HTTP.Host)
		require.NotEmpty(t, got.HTTP.Port)
	})
}
