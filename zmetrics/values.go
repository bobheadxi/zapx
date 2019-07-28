package zmetrics

import (
	"math"
	"time"

	"go.uber.org/zap/zapcore"
)

// ValueFunc returns the value associated with a field
type ValueFunc func(zapcore.Field) float64

// IncrementValue returns 1 when a field is encountered
func IncrementValue(zapcore.Field) float64 { return 1 }

// NumericalValue returns the numerical value of the field. It supports most basic
// numeric types such as Duration, Floats, and signed/unsigned integers.
func NumericalValue(f zapcore.Field) float64 {
	switch f.Type {
	case zapcore.DurationType:
		return float64(time.Duration(f.Integer))
	case zapcore.Float64Type:
		return float64(math.Float64frombits(uint64(f.Integer)))
	case zapcore.Float32Type:
		return float64(math.Float32frombits(uint32(f.Integer)))
	default:
		return float64(f.Integer)
	}
}

// UnixTimeValue returns the unix value of a time field.
func UnixTimeValue(f zapcore.Field) float64 {
	if f.Interface != nil {
		return float64(time.Unix(0, f.Integer).In(f.Interface.(*time.Location)).Unix())
	}
	// Fall back to UTC if location is nil.
	return float64(time.Unix(0, f.Integer).Unix())
}

// UnixNanoTimeValue returns the unix nanosecond value of a time field.
func UnixNanoTimeValue(f zapcore.Field) float64 {
	if f.Interface != nil {
		return float64(time.Unix(0, f.Integer).In(f.Interface.(*time.Location)).UnixNano())
	}
	// Fall back to UTC if location is nil.
	return float64(time.Unix(0, f.Integer).UnixNano())
}
