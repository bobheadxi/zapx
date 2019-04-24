package pgx

import (
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type databaseLogger struct {
	l    *zap.Logger
	opts Options
}

// Options denotes configuration for the pgx database logger
type Options struct {
	LogInfoAsInfo bool
}

// NewLogger wraps the given logger in pgx.Logger
func NewLogger(l *zap.Logger, opts Options) pgx.Logger {
	return &databaseLogger{
		// don't take stacktrace of wrapper class
		l.WithOptions(zap.AddCallerSkip(1)),
		opts,
	}
}

func (d *databaseLogger) Log(lv pgx.LogLevel, msg string, context map[string]interface{}) {
	var (
		ctxField   = zap.Any("pgx.context", context)
		levelField = zap.String("pgx.level", lv.String())
	)
	switch lv {
	case pgx.LogLevelDebug, pgx.LogLevelTrace:
		d.l.Debug(msg, levelField, ctxField)
	case pgx.LogLevelInfo:
		if d.opts.LogInfoAsInfo {
			d.l.Info(msg, levelField, ctxField)
		} else {
			d.l.Debug(msg, levelField, ctxField)
		}
	case pgx.LogLevelWarn:
		d.l.Warn(msg, levelField, ctxField)
	case pgx.LogLevelError:
		d.l.Error(msg, levelField, ctxField)
	}
}
