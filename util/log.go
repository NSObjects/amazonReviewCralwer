package util

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func init() {

	dyn := zap.NewAtomicLevel()
	dyn.SetLevel(zap.DebugLevel)
	cfg := zap.Config{
		Level:       dyn,
		Development: true,
		Encoding:    "console",

		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "logger",
			CallerKey:      "C",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},

		OutputPaths:      []string{"stdout", "./err.log"},
		ErrorOutputPaths: []string{"stderr"},
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger{
		l: l,
	}
}

type logger struct {
	l *zap.Logger
}

var Logger logger

func (l *logger) Info(msg string, fields ...zapcore.Field) {
	defer Logger.l.Sync()
	Logger.l.Info(msg, fields...)

}

func (l *logger) Error(msg string, fields ...zapcore.Field) {
	defer Logger.l.Sync()
	Logger.l.Error(msg, fields...)
}

func (l *logger) Debug(msg string, fields ...zapcore.Field) {
	defer Logger.l.Sync()
	Logger.l.Debug(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zapcore.Field) {
	defer Logger.l.Sync()
	Logger.l.Warn(msg, fields...)
}

func (l *logger) Panic(msg string, fields ...zapcore.Field) {
	defer Logger.l.Sync()
	Logger.l.Panic(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zapcore.Field) {
	defer Logger.l.Sync()
	Logger.l.Fatal(msg, fields...)
}
