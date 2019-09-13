package zapx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"go.bobheadxi.dev/zapx/ztest"
)

func TestWrapWithLevel(t *testing.T) {
	log, observed := ztest.NewObservable()

	log = log.WithOptions(WrapWithLevel(zap.ErrorLevel))
	log.Info("should not be logged")
	assert.Empty(t, observed.All())

	log.Error("should be logged")
	ztest.AssertObserved(t, []ztest.ObservedEntry{{Message: "should be logged"}}, observed)
}
