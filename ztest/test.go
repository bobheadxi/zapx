package ztest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// NewObservable bootstraps a logger that allows interrogation of output
func NewObservable() (*zap.Logger, *observer.ObservedLogs) {
	observer, out := observer.New(zap.InfoLevel)
	return zap.New(observer), out
}

// ObservedEntry is used for providing assertions on observed logs with AssertObserved()
type ObservedEntry struct {
	Message string
	Fields  map[string]interface{}
}

// AssertObserved asserts if the entries given were collected in order. It allows
// skips (entries in between the expected entries). Messages are asserted based
// on strings.Contains (not exact match)
func AssertObserved(t *testing.T, expected []ObservedEntry, observed *observer.ObservedLogs) bool {
	var i int
	for _, o := range observed.All() {
		if i > len(expected) {
			return true // all entries found
		}

		e := expected[i]
		if strings.Contains(o.Message, e.Message) {
			i++

			// map of discovered fields
			expectedFields := make(map[string]bool)
			for expectedKey := range e.Fields {
				expectedFields[expectedKey] = false
			}

			// check that each field is present and equal
			for _, field := range o.Context {
				if expectedValue, ok := e.Fields[field.Key]; ok {
					AssertEqualFields(t, zap.Any(field.Key, expectedValue), field)
					expectedFields[field.Key] = true
				}
			}

			// report any missed fields
			for k, v := range expectedFields {
				if !v {
					assert.Failf(t, "did not find expected field",
						"expected field '%s' but did not find it in message '%s'", k, e.Message)
				}
			}
		}
	}

	if i > len(expected) {
		return assert.Failf(t, "did not find expected message",
			"missed:\n%#v", expected[i:])
	}

	return !t.Failed()
}

// AssertEqualFields asserts that two zap.Fields are equal
func AssertEqualFields(t *testing.T, expected zap.Field, actual zap.Field) bool {
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Key, actual.Key)

	switch expected.Type {
	case zapcore.BinaryType, zapcore.ByteStringType:
		return assert.Equal(t, expected.Interface.([]byte), actual.Interface.([]byte))

	case zapcore.ArrayMarshalerType, zapcore.ObjectMarshalerType, zapcore.ErrorType, zapcore.ReflectType:
		return assert.EqualValues(t, expected.Interface, actual.Interface)

	default:
		return assert.EqualValues(t, expected, actual)
	}
}
