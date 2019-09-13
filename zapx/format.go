package zapx

import (
	"fmt"

	"go.uber.org/zap"
)

// FormatLogger exposes a common format-based logger interface
type FormatLogger interface {
	Errorf(string, ...interface{})
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

// ZapFormatLogger implements FormatLogger using the given zap logger
type ZapFormatLogger struct {
	L *zap.Logger
}

// NewFormatLogger instantiates a new default ZapFormatLogger
func NewFormatLogger(l *zap.Logger) ZapFormatLogger {
	return ZapFormatLogger{L: l}
}

// Errorf logs at the error level
func (z *ZapFormatLogger) Errorf(f string, v ...interface{}) {
	z.L.Error(fmt.Sprintf(f, v...))
}

// Warningf logs at the warn level
func (z *ZapFormatLogger) Warningf(f string, v ...interface{}) {
	z.L.Warn(fmt.Sprintf(f, v...))
}

// Infof logs at the info level
func (z *ZapFormatLogger) Infof(f string, v ...interface{}) {
	z.L.Info(fmt.Sprintf(f, v...))
}

// Debugf logs at the debug level
func (z *ZapFormatLogger) Debugf(f string, v ...interface{}) {
	z.L.Debug(fmt.Sprintf(f, v...))
}
