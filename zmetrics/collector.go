package zmetrics

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HistogramCollector defines a collector for histograms
type HistogramCollector struct {
	Histogram
	Value ValueFunc
}

// CounterCollector defines a collector for counters
type CounterCollector struct {
	Counter
	Value ValueFunc
}

// GaugeCollector defines a collector for gauges
type GaugeCollector struct {
	Gauge
	Value ValueFunc
}

// Collectors is a container for all your metrics collectors
type Collectors struct {
	Histograms map[string]HistogramCollector
	Counters   map[string]CounterCollector
	Gauges     map[string]GaugeCollector
}

// WithCollectors attaches collectors to the logger
func WithCollectors(l *zap.Logger, collectors Collectors) *zap.Logger {
	if collectors.Histograms == nil {
		collectors.Histograms = make(map[string]HistogramCollector)
	}
	if collectors.Counters == nil {
		collectors.Counters = make(map[string]CounterCollector)
	}
	if collectors.Gauges == nil {
		collectors.Gauges = make(map[string]GaugeCollector)
	}
	return zap.New(&collectorCore{l.Core(), collectors})
}

type collectorCore struct {
	c          zapcore.Core
	collectors Collectors
}

func (c *collectorCore) capture(fields []zapcore.Field) {
	for _, f := range fields {
		if h, ok := c.collectors.Histograms[f.Key]; ok {
			h.Observe(h.Value(f))
		}
		if c, ok := c.collectors.Counters[f.Key]; ok {
			c.Add(c.Value(f))
		}
		if v, ok := c.collectors.Gauges[f.Key]; ok {
			v.Set(v.Value(f))
		}
	}
}

func (c *collectorCore) With(fields []zapcore.Field) zapcore.Core {
	c.capture(fields)
	return c.c.With(fields)
}

func (c *collectorCore) Write(e zapcore.Entry, fields []zapcore.Field) error {
	c.capture(fields)
	return c.c.Write(e, fields)
}

func (c *collectorCore) Enabled(zapcore.Level) bool { return true }
func (c *collectorCore) Sync() error                { return c.c.Sync() }

func (c *collectorCore) Check(e zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return c.c.Check(e, checked)
}
