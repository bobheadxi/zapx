package metrics

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./stub/histogram.go --fake-name StubHistogram . Histogram
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./stub/counter.go --fake-name StubCounter . Counter
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./stub/gauge.go --fake-name StubGauge . Gauge

// Histogram describes a metric that takes repeated observations of the same
// kind of thing, and produces a statistical summary of those observations,
// typically expressed as quantiles or buckets.
type Histogram interface{ Observe(value float64) }

// Counter describes a metric that accumulates values monotonically.
type Counter interface{ Add(delta float64) }

// Gauge describes a metric that takes specific values over time.
type Gauge interface{ Set(value float64) }
