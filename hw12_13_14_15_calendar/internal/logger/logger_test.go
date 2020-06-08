package logger //nolint:golint,stylecheck

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Logger(t *testing.T) {
	t.Run("logging", func(t *testing.T) {
		filename := "temp.log"

		s := Settings{filename, ""}
		err := Configure(s)
		if err != nil {
			fmt.Println(err)
		}

		l := Logger
		l.Info("logger construction succeeded")

		file, err := ioutil.ReadFile(filename)
		require.NoError(t, err)
		got := string(file)
		err = os.Remove(filename)
		require.NoError(t, err)

		require.Contains(t, got, "logger construction succeeded\n")
	})
}
