package zapx

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a default zap logger based with basic defaults.
// If `logpath` is set, will use json encoder and write to the provided filepath.
// If `dev` is true, will use zap.DevelopmentConfig with colors enabled, though
// if logpath is set, the json encoder (without colors) will be enforced.
// If `dev` is false, will use zap.ProductionConfig.
func New(logpath string, dev bool) (l *zap.Logger, err error) {
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

	// instantiate logger
	if l, err = config.Build(); err != nil {
		return nil, fmt.Errorf("new logger: %s", err.Error())
	}

	return
}
