// platform/pkg/observability/logger/zap/logger.go
package zap

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/0xsj/scout/platform/pkg/observability/logger"
)

// zapLogger wraps zap.Logger to implement our Logger interface.
type zapLogger struct {
	zap *zap.Logger
}

// New creates a new production logger using Zap.
// Uses JSON encoding for structured logging.
func New(level logger.Level) (logger.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(toZapLevel(level))
	
	// JSON encoding for production
	config.Encoding = "json"
	
	// ISO8601 timestamps
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	
	// Use "level" instead of "severity"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	
	// Use "message" instead of "msg"
	config.EncoderConfig.MessageKey = "message"
	
	// Caller info
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	
	z, err := config.Build(
		zap.AddCallerSkip(1), // Skip wrapper functions
	)
	if err != nil {
		return nil, err
	}
	
	return &zapLogger{zap: z}, nil
}

// NewDevelopment creates a new development logger using Zap.
// Uses console encoding with colors for human readability.
func NewDevelopment(level logger.Level) (logger.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(toZapLevel(level))
	
	// Console encoding for development
	config.Encoding = "console"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	
	z, err := config.Build(
		zap.AddCallerSkip(1), // Skip wrapper functions
	)
	if err != nil {
		return nil, err
	}
	
	return &zapLogger{zap: z}, nil
}

// FromZap wraps an existing zap.Logger.
// Useful when you already have a configured zap logger.
func FromZap(z *zap.Logger) logger.Logger {
	return &zapLogger{zap: z}
}

func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.zap.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) With(keysAndValues ...interface{}) logger.Logger {
	newZap := l.zap.Sugar().With(keysAndValues...).Desugar()
	return &zapLogger{zap: newZap}
}

func (l *zapLogger) WithContext(ctx context.Context) logger.Logger {
	// Extract common context values
	// You can add custom extractors here for request ID, trace ID, etc.
	
	// For now, just return self
	// We'll implement context extraction when we build context propagation
	return l
}

func (l *zapLogger) WithError(err error) logger.Logger {
	return l.With("error", err.Error())
}

// Sync flushes any buffered log entries.
// Should be called before application shutdown.
func (l *zapLogger) Sync() error {
	return l.zap.Sync()
}

// toZapLevel converts our Level to zap's zapcore.Level
func toZapLevel(level logger.Level) zapcore.Level {
	switch level {
	case logger.DebugLevel:
		return zapcore.DebugLevel
	case logger.InfoLevel:
		return zapcore.InfoLevel
	case logger.WarnLevel:
		return zapcore.WarnLevel
	case logger.ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}