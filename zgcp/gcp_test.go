package zgcp

import (
	"net/http"
	"testing"

	"cloud.google.com/go/errorreporting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bobheadxi.dev/zapx/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type fakeGCPReporter struct {
	lastEntry errorreporting.Entry
	flushed   bool
}

func (f *fakeGCPReporter) Report(e errorreporting.Entry) {
	f.lastEntry = errorreporting.Entry{
		Error: e.Error,
		Req:   e.Req,
		User:  e.User,
		Stack: []byte(string(e.Stack)),
	}
}
func (f *fakeGCPReporter) Flush() { f.flushed = true }

func Test_gcpErrorReportingZapCore_Enabled(t *testing.T) {
	var f fakeGCPReporter
	core := gcpErrorsWrapCore(&f, Fields{})
	assert.False(t, core.Enabled(zapcore.DebugLevel))
	assert.False(t, core.Enabled(zapcore.InfoLevel))
	assert.True(t, core.Enabled(zapcore.WarnLevel))
	assert.True(t, core.Enabled(zapcore.ErrorLevel))
}

func Test_stacktrace(t *testing.T) {
	buf := internal.NewBuffer().(*internal.Buffer)
	stack := stacktrace(buf.Bytes())
	assert.NotNil(t, stack)
}

func Test_gcpErrorReportingZapCore_Write(t *testing.T) {
	type args struct {
		entry  zapcore.Entry
		fields []zapcore.Field
	}
	type want struct {
		user          string
		requestMethod string
		err           string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{"all empty", args{
			zapcore.Entry{}, nil,
		}, want{}, false},
		{"message", args{
			zapcore.Entry{Message: "hello world"}, nil,
		}, want{err: "hello world"}, false},
		{"request ID", args{
			zapcore.Entry{},
			[]zapcore.Field{
				zap.String("request-id", "1234"),
			},
		}, want{user: "1234"}, false},
		{"http.Request", args{
			zapcore.Entry{},
			[]zapcore.Field{
				zap.Any("request", &http.Request{Method: "POST"}),
			},
		}, want{requestMethod: "POST"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f fakeGCPReporter
			core := gcpErrorsWrapCore(&f, Fields{
				UserKey:    "request-id",
				RequestKey: "request",
			})
			if tt.wantErr {
				assert.Error(t, core.Write(tt.args.entry, tt.args.fields))
			} else {
				assert.NoError(t, core.Write(tt.args.entry, tt.args.fields))
			}
			assert.Equal(t, tt.want.user, f.lastEntry.User)
			if tt.want.requestMethod != "" {
				require.NotNil(t, f.lastEntry.Req)
				assert.Equal(t, tt.want.requestMethod, f.lastEntry.Req.Method)
			}
			assert.Contains(t, f.lastEntry.Error.Error(), tt.want.err)
		})
	}
}

func Test_gcpErrorReportingZapCore_Sync(t *testing.T) {
	var f fakeGCPReporter
	core := gcpErrorsWrapCore(&f, Fields{})
	assert.False(t, f.flushed)
	assert.NoError(t, core.Sync())
	assert.True(t, f.flushed)
}
