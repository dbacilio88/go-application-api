package log

import (
	"context"
	"fmt"
	"github.com/dbacilio88/go-application-api/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Core zapcore.Core
var LoggerInstance *zap.Logger
var Ctx context.Context

type RequestIdKey struct{}
type ctxKey struct {
}

func LoggerConfiguration(logLevel string) (*zap.Logger, error) {

	var level zap.AtomicLevel

	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)

	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	fileEncoder := zapcore.NewJSONEncoder(zapEncoderConfig)

	logFile := fmt.Sprintf("/logs/%s-%s.log", config.Configuration.Microservices.Name, config.Configuration.Environment.Value)

	logRotation := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // max size in megabytes before log rotation
		MaxBackups: 10, // max number of old log files to keep
		MaxAge:     30, // max number of days to retain old log file
		LocalTime:  true,
		Compress:   true, // compress old log file
	}

	defaultLogLevel := level

	Core = zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logRotation), defaultLogLevel).With(
			[]zap.Field{
				zap.String("app", config.Configuration.Microservices.Name),
				zap.String("env", config.Configuration.Environment.Value),
				zap.String("log", config.Configuration.Log.Level),
			}),
		zapcore.NewCore(fileEncoder, zapcore.Lock(os.Stdout), defaultLogLevel).With(
			[]zap.Field{
				zap.String("app", config.Configuration.Microservices.Name),
				zap.String("env", config.Configuration.Environment.Value),
				zap.String("log", config.Configuration.Log.Level),
			}),
	)

	LoggerInstance = zap.New(Core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return LoggerInstance, nil
}

func Reset() {
	Ctx = nil
	clone := Core.With([]zap.Field{})
	LoggerInstance = zap.New(clone)
}

func FromCtx() *zap.Logger {
	if lp, ok := Ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return lp
	} else if lp := LoggerInstance; lp != nil {
		return lp
	}
	return zap.NewNop()
}

func WithCtx(ctx context.Context, logger *zap.Logger) context.Context {
	Ctx = ctx
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == logger {
			return ctx
		}
	}
	return context.WithValue(ctx, ctxKey{}, logger)
}

func GetRequestId() zap.Field {
	return zap.Any("requestId", Ctx.Value("requestId"))
}
func GetRequestIdValue() string {
	value := Ctx.Value("requestId")

	if id, ok := value.(uuid.UUID); ok {
		return id.String()
	}
	return ""
}
