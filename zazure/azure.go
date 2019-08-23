package zazure

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ServiceConfig defines configuration for an Azure service
type ServiceConfig struct {
}

// NewAzureLogger https://docs.microsoft.com/en-us/azure/azure-monitor/overview
//
// https://docs.microsoft.com/en-us/azure/azure-monitor/platform/runbook-datacollect
// https://docs.microsoft.com/en-us/azure/azure-monitor/platform/data-collector-api
func NewAzureLogger(
	l *zap.Logger,
	service ServiceConfig,
	debug bool,
) (*zap.Logger, error) {
	l.Info("setting up azure logs sink")

	return l.
		WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, azureWrapCore())
		})), nil
}

type azureZapCore struct {
	enc zapcore.Encoder
}

func azureWrapCore() zapcore.Core {
	return &azureZapCore{
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			NameKey:    "logger",
			MessageKey: "msg",
			LevelKey:   "level",

			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
	}
}

func (z *azureZapCore) Enabled(l zapcore.Level) bool {
	return l >= zapcore.WarnLevel
}

func (z *azureZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// encode everything as the message
	buf, err := z.enc.EncodeEntry(entry, fields)
	if err != nil {
		return fmt.Errorf("failed to encode entry for gcp error reporter: %v", err)
	}

	// TODO
	println(buf)

	return nil
}

func (z *azureZapCore) Sync() error {
	return nil
}

func (z *azureZapCore) With(fields []zapcore.Field) zapcore.Core {
	clone := &azureZapCore{z.enc.Clone()}
	for i := range fields {
		fields[i].AddTo(clone.enc)
	}
	return clone
}

func (z *azureZapCore) Check(e zapcore.Entry, c *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(e.Level) {
		return c.AddCore(e, z)
	}
	return c
}
