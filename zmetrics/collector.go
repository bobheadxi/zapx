package zmetrics

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.bobheadxi.dev/zapx/zmetrics/metrics"
)

// HistogramCollector defines a collector for histograms
type HistogramCollector struct {
	metrics.Histogram
	Value ValueFunc
}

// NewHistogram instantiates a new histogram collector
func NewHistogram(m metrics.Histogram, v ValueFunc) HistogramCollector {
	return HistogramCollector{m, v}
}

// CounterCollector defines a collector for counters
type CounterCollector struct {
	metrics.Counter
	Value ValueFunc
}

// NewCounter instantiates a new counter collector
func NewCounter(m metrics.Counter, v ValueFunc) CounterCollector {
	return CounterCollector{m, v}
}

// GaugeCollector defines a collector for gauges
type GaugeCollector struct {
	metrics.Gauge
	Value ValueFunc
}

// NewGauge instantiates a new gauge collector
func NewGauge(m metrics.Gauge, v ValueFunc) GaugeCollector {
	return GaugeCollector{m, v}
}

// Collectors is a container for all your metrics collectors. It should not be
// modified after creation.
type Collectors struct {
	Histograms map[string]HistogramCollector
	Counters   map[string]CounterCollector
	Gauges     map[string]GaugeCollector
}

// ExactMatch tries to match the given field exactly
func (c *Collectors) ExactMatch(f zapcore.Field) {
	if h, ok := c.Histograms[f.Key]; ok {
		h.Observe(h.Value(f))
	}
	if c, ok := c.Counters[f.Key]; ok {
		c.Add(c.Value(f))
	}
	if v, ok := c.Gauges[f.Key]; ok {
		v.Set(v.Value(f))
	}
}

// Collect attaches collectors to the logger
func Collect(l *zap.Logger, c *Collectors) *zap.Logger {
	if c.Histograms == nil {
		c.Histograms = make(map[string]HistogramCollector)
	}
	if c.Counters == nil {
		c.Counters = make(map[string]CounterCollector)
	}
	if c.Gauges == nil {
		c.Gauges = make(map[string]GaugeCollector)
	}
	return zap.New(wrap(l.Core(), c))
}
