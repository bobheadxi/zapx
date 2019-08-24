/*

Package bench is just a bunch of benchmarks for the various subpackages in `zapx`.

Learn more about benchmarks here: https://golang.org/pkg/testing/#hdr-Benchmarks

To run these benchmarks:

	go test -bench . -benchmem ./...

Results can be exported using https://github.com/bobheadxi/gobenchdata - for
example, to export a run as JSON:

	go test -bench . -benchmem ./... | gobenchdata --json bench.json

*/
package benchmarks

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"go.bobheadxi.dev/zapx/internal/testdata"
	"go.bobheadxi.dev/zapx/zapx"
	"go.bobheadxi.dev/zapx/zgcp"
)

func BenchmarkLoggerInfo(b *testing.B) {
	os.RemoveAll("./tmp")
	defer os.RemoveAll("./tmp")

	b.Run("zapx::plain", func(b *testing.B) {
		logger, err := zapx.New("./tmp/"+b.Name(), false, zapx.OnlyToFile())
		require.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(testdata.FakeMessage(), testdata.FakeFields()...)
			}
		})
		logger.Sync()
	})
	b.Run("zapx/zgcp::error-reporting", func(b *testing.B) {
		srv := newTestServer()
		defer srv.Close()

		logger, err := zapx.New("./tmp/"+b.Name(), false, zapx.OnlyToFile())
		require.NoError(b, err)
		logger, err = zgcp.NewErrorReportingLogger(logger, zgcp.ServiceConfig{}, zgcp.Fields{
			UserKey: "obj1",
		}, false, option.WithEndpoint(srv.URL), option.WithCredentials(&google.Credentials{}))
		require.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(testdata.FakeMessage(), testdata.FakeFields()...)
			}
		})
		logger.Sync()
	})
	/* TODO
	b.Run("zapx/zpgx::logger", func(b *testing.B) {})
	*/
}

func BenchmarkLoggerError(b *testing.B) {
	os.RemoveAll("./tmp")
	defer os.RemoveAll("./tmp")

	b.Run("zapx::plain", func(b *testing.B) {
		logger, err := zapx.New("./tmp/"+b.Name(), false, zapx.OnlyToFile())
		require.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Error(testdata.FakeMessage(), testdata.FakeFields()...)
			}
		})
		logger.Sync()
	})
	b.Run("zapx/zgcp::error-reporting", func(b *testing.B) {
		srv := newTestServer()
		defer srv.Close()

		logger, err := zapx.New("./tmp/"+b.Name(), false, zapx.OnlyToFile())
		require.NoError(b, err)
		logger, err = zgcp.NewErrorReportingLogger(logger, zgcp.ServiceConfig{}, zgcp.Fields{
			UserKey: "obj1",
		}, false, option.WithEndpoint(srv.URL), option.WithCredentials(&google.Credentials{}))
		require.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Error(testdata.FakeMessage(), testdata.FakeFields()...)
			}
		})
		logger.Sync()
	})
	/* TODO
	b.Run("zapx/zpgx::logger", func(b *testing.B) {})
	*/
}

/* TODO
func BenchmarkMiddleware(b *testing.B) {

}
*/
