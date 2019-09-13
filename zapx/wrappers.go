package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WrapWithLevel is a zap.Option that spawns a child logger at the specified level
func WrapWithLevel(level zapcore.Level) zap.Option {
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &coreWithLevel{Core: core, level: level}
	})
}

type coreWithLevel struct {
	zapcore.Core
	level zapcore.Level
}

func (c *coreWithLevel) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

func (c *coreWithLevel) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !c.level.Enabled(e.Level) {
		return ce
	}
	return c.Core.Check(e, ce)
}
