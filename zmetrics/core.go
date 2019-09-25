package zmetrics

import (
	"sync"

	"go.uber.org/zap/zapcore"
)

type collectorCore struct {
	wrapped    zapcore.Core
	collectors *Collectors

	withFields []zapcore.Field
	fieldMux   sync.RWMutex
}

func wrap(wrapped zapcore.Core, c *Collectors) zapcore.Core {
	return &collectorCore{wrapped: wrapped, collectors: c, withFields: make([]zapcore.Field, 0)}
}

func (c *collectorCore) capture(fields []zapcore.Field) {
	for _, f := range fields {
		c.collectors.ExactMatch(f)
	}
	c.fieldMux.RLock()
	for _, f := range c.withFields {
		c.collectors.ExactMatch(f)
	}
	c.fieldMux.RUnlock()
}

func (c *collectorCore) With(fields []zapcore.Field) zapcore.Core {
	c.fieldMux.Lock()
	withFields := append(c.withFields, fields...)
	c.fieldMux.Unlock()
	return &collectorCore{
		wrapped:    c.wrapped.With(fields),
		collectors: c.collectors,
		withFields: withFields,
	}
}

func (c *collectorCore) Write(e zapcore.Entry, fields []zapcore.Field) error {
	c.capture(fields)

	return c.wrapped.Write(e, fields)
}

func (c *collectorCore) Enabled(l zapcore.Level) bool { return c.wrapped.Enabled(l) }

func (c *collectorCore) Sync() error {
	return c.wrapped.Sync()
}

func (c *collectorCore) Check(e zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(e.Level) {
		return checked.AddCore(e, c)
	}
	return c.wrapped.Check(e, checked)
}
