package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_fileInPath(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		err := fileInPath("./main_test.go")
		require.Nil(t, err)
	})
	t.Run("not existing", func(t *testing.T) {
		err := fileInPath("./file.net")
		require.NotNil(t, err)
	})
	t.Run("not regular", func(t *testing.T) {
		err := fileInPath("../go-validate")
		require.Equal(t, err, ErrFileNotRegular)
	})
}
