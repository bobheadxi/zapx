package zapx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type args struct {
		logpath string
		dev     bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"dev-no-path", args{"", true}, false},
		{"prod-no-path", args{"", false}, false},
		{"dev-with-path", args{"./tmp/log", true}, false},
		{"prod-with-path", args{"./tmp/log", false}, false},
		{"bad-dir-dev", args{"/root/toor", true}, true},
		{"bad-dir-prod", args{"/root/toor", false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.logpath, tt.args.dev)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
