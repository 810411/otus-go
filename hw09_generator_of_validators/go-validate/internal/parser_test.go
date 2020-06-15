package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getValidatorStr(t *testing.T) {
	t.Run("empty tag", func(t *testing.T) {
		v := getValidatorStr("")
		require.Empty(t, v)
	})
	t.Run("validation tag", func(t *testing.T) {
		v := getValidatorStr(`validate:"len:5"`)
		require.Equal(t, v, "len:5")
	})
	t.Run("ignoring tag", func(t *testing.T) {
		v := getValidatorStr(`json:"omitempty"`)
		require.Empty(t, v)
	})
	t.Run("mixed tag", func(t *testing.T) {
		v := getValidatorStr(`json:"id" validate:"len:36"`)
		require.Equal(t, v, "len:36")
	})
}
