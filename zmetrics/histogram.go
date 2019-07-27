package zmetrics

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WithHistogram sets an histogram on the given logger.
func WithHistogram(l *zap.Logger, key string, value ValueFunc, hst Histogram) *zap.Logger {
	return zap.New(&histogramWriter{l.Core(), key, value, hst})
}

type histogramWriter struct {
	c zapcore.Core

	key   string
	value ValueFunc
	hst   Histogram
}

func (w *histogramWriter) capture(fields []zapcore.Field) {
	for _, f := range fields {
		if f.Key == w.key {
			w.hst.Observe(w.value(f))
		}
	}
}

func (w *histogramWriter) With(fields []zapcore.Field) zapcore.Core {
	w.capture(fields)
	return w.c.With(fields)
}

func (w *histogramWriter) Write(e zapcore.Entry, fields []zapcore.Field) error {
	w.capture(fields)
	return w.c.Write(e, fields)
}

func (w *histogramWriter) Enabled(zapcore.Level) bool { return true }
func (w *histogramWriter) Sync() error                { return w.c.Sync() }

func (w *histogramWriter) Check(e zapcore.Entry, c *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return w.c.Check(e, c)
}
