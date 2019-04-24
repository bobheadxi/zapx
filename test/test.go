package test

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// NewObservable bootstraps a logger that allows interrogation of output
func NewObservable() (*zap.Logger, *observer.ObservedLogs) {
	observer, out := observer.New(zap.InfoLevel)
	return zap.New(observer), out
}
