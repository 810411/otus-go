package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Settings struct {
	File  string
	Level string
}

var Logger *zap.Logger

func Configure(s Settings) (err error) {
	var level zap.AtomicLevel
	switch s.Level {
	case "warn":
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	out := []string{"stdout"}
	if s.File != "" {
		out = append(out, s.File)
	}

	cfg := zap.Config{
		Level:       level,
		Encoding:    "console",
		OutputPaths: out,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
		},
	}

	Logger, err = cfg.Build()
	if err != nil {
		return
	}

	_ = Logger.Sync()

	return
}
