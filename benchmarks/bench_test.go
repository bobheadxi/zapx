/*

Package bench is just a bunch of benchmarks for the various subpackages in `zapx`.

*/
package benchmarks

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/bobheadxi/zapx"
	"github.com/bobheadxi/zapx/zgcp"
)

func BenchmarkLoggingInfo(b *testing.B) {
	os.RemoveAll("./tmp")
	defer os.RemoveAll("./tmp")

	b.Run("zapx::plain-logger", func(b *testing.B) {
		logger, err := zapx.New("./tmp/"+b.Name(), false, zapx.OnlyToFile())
		require.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(fakeMessage(), fakeFields()...)
			}
		})
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
				logger.Info(fakeMessage(), fakeFields()...)
			}
		})
	})
	/* TODO
	b.Run("zapx/zpgx::logger", func(b *testing.B) {})
	*/
}

func BenchmarkLoggingError(b *testing.B) {
	os.RemoveAll("./tmp")
	defer os.RemoveAll("./tmp")

	b.Run("zapx::plain-logger", func(b *testing.B) {
		logger, err := zapx.New("./tmp/"+b.Name(), false, zapx.OnlyToFile())
		require.NoError(b, err)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Error(fakeMessage(), fakeFields()...)
			}
		})
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
				logger.Error(fakeMessage(), fakeFields()...)
			}
		})
	})
	/* TODO
	b.Run("zapx/zpgx::logger", func(b *testing.B) {})
	*/
}

/* TODO
func BenchmarkMiddleware(b *testing.B) {

}
*/
