package zmetrics

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WithGauge sets a gauge on the given logger.
func WithGauge(l *zap.Logger, key string, value ValueFunc, gg Gauge) *zap.Logger {
	return zap.New(&gaugeWriter{l.Core(), key, value, gg})
}

type gaugeWriter struct {
	c zapcore.Core

	key   string
	value ValueFunc
	gg    Gauge
}

func (w *gaugeWriter) capture(fields []zapcore.Field) {
	for _, f := range fields {
		if f.Key == w.key {
			w.gg.Set(w.value(f))
		}
	}
}

func (w *gaugeWriter) With(fields []zapcore.Field) zapcore.Core {
	w.capture(fields)
	return w.c.With(fields)
}

func (w *gaugeWriter) Write(e zapcore.Entry, fields []zapcore.Field) error {
	w.capture(fields)
	return w.c.Write(e, fields)
}

func (w *gaugeWriter) Enabled(zapcore.Level) bool { return true }
func (w *gaugeWriter) Sync() error                { return w.c.Sync() }

func (w *gaugeWriter) Check(e zapcore.Entry, c *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return w.c.Check(e, c)
}
