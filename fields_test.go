package zapx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestFields(t *testing.T) {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig)
	buf, err := encoder.EncodeEntry(
		zapcore.Entry{},
		[]zap.Field{
			Fields("hello",
				zap.String("bob", "head"),
				zap.Int("cool", 100)),
		})
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"hello":{"bob":"head","cool":100}`)
}
