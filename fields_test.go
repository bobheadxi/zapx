package zapx

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bobheadxi/zapx/benchmarks/testdata"
)

func TestFieldSet(t *testing.T) {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig)
	buf, err := encoder.EncodeEntry(
		zapcore.Entry{},
		[]zap.Field{
			FieldSet("hello",
				zap.String("bob", "head"),
				zap.Int("cool", 100)),
		})
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), `"hello":{"bob":"head","cool":100}`)
}

func BenchmarkFieldSet(b *testing.B) {
	os.RemoveAll("./tmp")
	defer os.RemoveAll("./tmp")

	b.Run("zapx.FieldSet(...zap.Field)", func(b *testing.B) {
		logger, err := New("./tmp/"+b.Name(), false, OnlyToFile())
		require.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(testdata.FakeMessage(),
					FieldSet("data",
						zap.Int("ID", 42),
						zap.String("Name", "bobheadxi"),
						zap.Time("CreatedAt", time.Date(1998, 3, 11, 12, 0, 0, 0, time.UTC))))
			}
		})
		logger.Sync()
	})
	b.Run("zap.Any(map[string]interface{})", func(b *testing.B) {
		logger, err := New("./tmp/"+b.Name(), false, OnlyToFile())
		require.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(testdata.FakeMessage(),
					zap.Any("data", map[string]interface{}{
						"ID":        42,
						"Name":      "bobheadxi",
						"CreatedAt": time.Date(1998, 3, 11, 12, 0, 0, 0, time.UTC)}))
			}
		})
		logger.Sync()
	})
	b.Run("zap.Any(zapcore::ObjectMarshaller)", func(b *testing.B) {
		logger, err := New("./tmp/"+b.Name(), false, OnlyToFile())
		require.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(testdata.FakeMessage(),
					zap.Object("data", testdata.FakeObject()))
			}
		})
		logger.Sync()
	})
}
