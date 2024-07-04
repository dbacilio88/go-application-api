package log

import (
	"fmt"
	"github.com/dbacilio88/go-application-api/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Core zapcore.Core
var LoggerInstance *zap.Logger

func ApplyLoggerConfiguration(logLevel string) (*zap.Logger, error) {
	var level zap.AtomicLevel

	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:  "ts",
		LevelKey: "level",
		NameKey:  "logger",
	}

	fileEncoder := zapcore.NewJSONEncoder(zapEncoderConfig)

	logFile := fmt.Sprintf("/logs/%s-%s.log", config.Configuration.Microservices.Name, config.Configuration.Environment)

	logRotation := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   true,
	}

	defaultLogLevel := level

	Core = zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logRotation), defaultLogLevel).With(
			[]zap.Field{
				zap.String("app", config.Configuration.Microservices.Name),
				zap.String("env", config.Configuration.Microservices.Name),
			}),
		zapcore.NewCore(fileEncoder, zapcore.Lock(os.Stdout), defaultLogLevel).With(
			[]zap.Field{
				zap.String("app", config.Configuration.Microservices.Name),
				zap.String("env", config.Configuration.Microservices.Name),
			}),
	)

	LoggerInstance = zap.New(Core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return LoggerInstance, nil

}
