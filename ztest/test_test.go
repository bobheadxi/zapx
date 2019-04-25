package ztest

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewObservable(t *testing.T) {
	logger, out := NewObservable()
	logger.Info("hi")
	require.True(t, len(out.All()) > 0, "log should have been captured")
	assert.Equal(t, out.All()[0].Message, "hi")
	assert.Equal(t, out.All()[0].Level, zapcore.InfoLevel)
}
