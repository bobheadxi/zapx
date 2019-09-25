package zmetrics

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"gotest.tools/assert"

	"go.uber.org/zap/zaptest"

	"go.bobheadxi.dev/zapx/zmetrics/metrics/stub"
)

func TestWrappedLogger(t *testing.T) {
	var (
		hist  stub.StubHistogram
		count stub.StubCounter
		gauge stub.StubGauge
	)

	// set up logger
	l := zaptest.NewLogger(t)
	cl := Collect(l, &Collectors{
		Histograms: map[string]HistogramCollector{
			"duration": NewHistogram(&hist, NumericValue),
		},
		Counters: map[string]CounterCollector{
			"error": NewCounter(&count, IncrementValue),
		},
		Gauges: map[string]GaugeCollector{
			"queue_size": NewGauge(&gauge, NumericValue),
		},
	})

	cl.Info("some duration", zap.Duration("duration", 12*time.Second))
	assert.Equal(t, float64(12*time.Second), hist.ObserveArgsForCall(0))

	cl.Error("some error", zap.Error(errors.New("oh no")))
	assert.Equal(t, float64(1), count.AddArgsForCall(0))

	cl.Debug("status update", zap.Int("queue_size", 12))
	cl.Debug("status update 2", zap.Int("queue_size", 14))
	assert.Equal(t, float64(12), gauge.SetArgsForCall(0))
	assert.Equal(t, float64(14), gauge.SetArgsForCall(1))

	cl = cl.With(zap.Duration("duration", 3*time.Second))
	cl.Info("with field")
	//assert.Equal(t, float64(3*time.Second), hist.ObserveArgsForCall(1))
}
