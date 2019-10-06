package zapx

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewToWriter instantiates a simple zap logger that writes to the given writer
func NewToWriter(
	w io.Writer,
	level zapcore.Level,
	enc zapcore.Encoder,
) *zap.Logger {
	return zap.New(zapcore.NewCore(enc, zapcore.AddSync(w), zap.NewAtomicLevelAt(level)))
}
