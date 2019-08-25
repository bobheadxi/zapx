package zgcp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"net/http"
	"runtime"

	"cloud.google.com/go/errorreporting"
	"go.bobheadxi.dev/zapx/util/pool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/option"
)

// ServiceConfig defines configuration for a Google Cloud Platform service
type ServiceConfig struct {
	ProjectID string
	Name      string
	Version   string
}

// Fields defines special fields to parse out of a set of Zap fields. If empty,
// then will be ignored.
type Fields struct {
	UserKey    string
	RequestKey string
}

// NewErrorReportingLogger attaches the GCP Stackdriver Error Reporting client to
// the given zap logger. See https://cloud.google.com/error-reporting
func NewErrorReportingLogger(
	l *zap.Logger,
	service ServiceConfig,
	fields Fields,
	debug bool,
	opts ...option.ClientOption,
) (*zap.Logger, error) {
	l.Info("setting up GCP error reporting",
		zap.String("project_id", service.ProjectID))

	errHandler := func(error) {}
	if debug {
		errHandler = func(e error) { stdlog.Printf("gcp.error-reporter: %v\n", e) }
	}

	reporter, err := errorreporting.NewClient(
		context.Background(),
		service.ProjectID,
		errorreporting.Config{
			ServiceName:    service.Name,
			ServiceVersion: service.Version,
			OnError:        errHandler,
		},
		opts...)
	if err != nil {
		return nil, err
	}

	return l.
		WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, gcpErrorsWrapCore(reporter, fields))
		})), nil
}

type gcpReporter interface {
	Report(e errorreporting.Entry)
	Flush()
}

type gcpErrorReportingZapCore struct {
	reporter gcpReporter
	enc      zapcore.Encoder
	fields   Fields

	buffers pool.ByteBufferPool
}

func gcpErrorsWrapCore(reporter gcpReporter, fields Fields) zapcore.Core {
	return &gcpErrorReportingZapCore{
		reporter,
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
		fields,
		pool.NewByteBufferPool(),
	}
}

func (z *gcpErrorReportingZapCore) Enabled(l zapcore.Level) bool {
	return l >= zapcore.WarnLevel
}

func (z *gcpErrorReportingZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// extract relevant values from fields
	var user string
	var request *http.Request
	for _, f := range fields {
		if z.fields.UserKey != "" && f.Key == z.fields.UserKey && f.Type == zapcore.StringType {
			user = f.String
		}
		if z.fields.RequestKey != "" && f.Key == z.fields.RequestKey {
			request, _ = f.Interface.(*http.Request)
		}
	}

	// encode everything as the message
	buf, err := z.enc.EncodeEntry(entry, fields)
	if err != nil {
		return fmt.Errorf("failed to encode entry for gcp error reporter: %v", err)
	}

	// report to GCP
	stackBuf := z.buffers.Get()
	z.reporter.Report(errorreporting.Entry{
		Error: errors.New(buf.String()),
		User:  user,
		Req:   request,

		// GCP Error Reporting does not like Zap's custom stacktraces (from entry.Stack),
		// so a custom stacktrace must be taken that conforms to the standard. Ugh.
		// See stacktrace() for details.
		Stack: stacktrace(stackBuf.Bytes()),
	})
	stackBuf.Free()

	return nil
}

func (z *gcpErrorReportingZapCore) Sync() error {
	z.reporter.Flush()
	return nil
}

func (z *gcpErrorReportingZapCore) With(fields []zapcore.Field) zapcore.Core {
	clone := &gcpErrorReportingZapCore{z.reporter, z.enc.Clone(), z.fields, z.buffers}
	for i := range fields {
		fields[i].AddTo(clone.enc)
	}
	return clone
}

func (z *gcpErrorReportingZapCore) Check(e zapcore.Entry, c *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(e.Level) {
		return c.AddCore(e, z)
	}
	return c
}

// stackSkip denotes the number of stack levels we need to skip over:
// * stacktrace()
// * Core.Write()
// * CheckedEntry.Write()
// * logger.Warn() or equivalent
const stackSkip = 4

var lineSep = []byte{'\n'}

// stacktrace captures the calling goroutine's stack and trims out irrelevant
// levels (as described in documentation for stackSkip)
func stacktrace(buf []byte) []byte {
	runtime.Stack(buf, false)
	lines := bytes.Split(buf, lineSep)
	lines = append(lines[:1], lines[2*stackSkip+1:]...)
	return bytes.Join(lines, lineSep)
}
