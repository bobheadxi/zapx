package zapx

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// New creates a default zap logger based with basic defaults.
// If `logpath` is set, will use json encoder and write to the provided filepath.
// If `dev` is true, will use zap.DevelopmentConfig with colors enabled, though
// if logpath is set, the json encoder (without colors) will be enforced.
// If `dev` is false, will use zap.ProductionConfig.
func New(logpath string, dev bool, opts ...Option) (l *zap.Logger, err error) {
	var config zap.Config
	if dev {
		// Log:         DebugLevel
		// Encoder:     console
		// Errors:      stderr
		// Sampling:    no
		// Stacktraces: WarningLevel
		// Colors:      capitals
		config = zap.NewDevelopmentConfig()
		if logpath == "" {
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
	} else {
		// Log:         InfoLevel
		// Encoder:     json
		// Errors:      stderr
		// Sampling:    yes
		// Stacktraces: ErrorLevel
		config = zap.NewProductionConfig()
	}

	// set log output configuration if provided
	if logpath != "" {
		if err = os.MkdirAll(filepath.Dir(logpath), os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directories for logpath '%s': %s",
				logpath, err.Error())
		}
		config.OutputPaths = append(config.OutputPaths, logpath)
		config.Encoding = "json"
	}

	// apply additional options
	for _, apply := range opts {
		if err := apply(&config); err != nil {
			return nil, err
		}
	}

	// instantiate logger
	if l, err = config.Build(); err != nil {
		return nil, fmt.Errorf("new logger: %s", err.Error())
	}

	return
}
