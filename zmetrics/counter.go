package zmetrics

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WithCounter sets a counter on the given logger.
func WithCounter(l *zap.Logger, key string, value ValueFunc, gg Gauge) *zap.Logger {
	return zap.New(&gaugeWriter{l.Core(), key, value, gg})
}

type counterWriter struct {
	c zapcore.Core

	key   string
	value ValueFunc
	gg    Gauge
}

func (w *counterWriter) capture(fields []zapcore.Field) {
	for _, f := range fields {
		if f.Key == w.key {
			w.gg.Set(w.value(f))
		}
	}
}

func (w *counterWriter) With(fields []zapcore.Field) zapcore.Core {
	w.capture(fields)
	return w.c.With(fields)
}

func (w *counterWriter) Write(e zapcore.Entry, fields []zapcore.Field) error {
	w.capture(fields)
	return w.c.Write(e, fields)
}

func (w *counterWriter) Enabled(zapcore.Level) bool { return true }
func (w *counterWriter) Sync() error                { return w.c.Sync() }

func (w *counterWriter) Check(e zapcore.Entry, c *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return w.c.Check(e, c)
}
