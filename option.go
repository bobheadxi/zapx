package zapx

import (
	"errors"

	"go.uber.org/zap"
)

// Option denotes an option for the zapx constructor
type Option func(*zap.Config) error

// OnlyToFile removes all other outputs and only writes to file
func OnlyToFile() Option {
	return func(cfg *zap.Config) error {
		if len(cfg.OutputPaths) == 1 {
			return errors.New("no output file set")
		}
		cfg.OutputPaths = cfg.OutputPaths[1:]
		return nil
	}
}

// WithFields sets initial fields in the configuration
func WithFields(fields map[string]interface{}) Option {
	return func(cfg *zap.Config) error {
		cfg.InitialFields = fields
		return nil
	}
}

// WithDebug sets log level to debug if given bool is true
func WithDebug(debug bool) Option {
	return func(cfg *zap.Config) error {
		if debug {
			cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		}
		return nil
	}
}
