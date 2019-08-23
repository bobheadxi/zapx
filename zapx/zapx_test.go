package zapx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	os.RemoveAll("tmp")
	type args struct {
		logpath string
		dev     bool
		opts    []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ok: dev, no path", args{"", true, nil}, false},
		{"ok: prod, no path", args{"", false, nil}, false},
		{"ok: dev, with path", args{"./tmp/log", true, nil}, false},
		{"ok: prod, with path", args{"./tmp/log", false, nil}, false},
		{"ok: with fields", args{"./tmp/log", false, []Option{
			WithFields(map[string]interface{}{"hello": "world"}),
		}}, false},
		{"ok: with only file path", args{"./tmp/log", false, []Option{
			OnlyToFile(),
		}}, false},
		{"fail: only filepath option without file", args{"", false, []Option{
			OnlyToFile(),
		}}, true},
		{"fail: bad path", args{"/root/toor", true, nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.logpath, tt.args.dev, tt.args.opts...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
