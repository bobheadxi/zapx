package benchmarks

import (
	"errors"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type obj struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

func (o *obj) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", o.ID)
	enc.AddString("name", o.Name)
	enc.AddTime("created_at", o.CreatedAt)
	return nil
}

type objs []*obj

func (o objs) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	var err error
	for _, v := range o {
		err = multierr.Append(err, arr.AppendObject(v))
	}
	return err
}

var (
	errExample  = errors.New("hello world")
	dataInts    = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	dataStrings = []string{"the", "quick", "brown", "fox", "jumped", "over", "the", "sleeping", "dog", "amazing"}
	dataTimes   = []time.Time{
		time.Unix(0, 0),
		time.Unix(1, 0),
		time.Unix(2, 0),
		time.Unix(3, 0),
		time.Unix(4, 0),
		time.Unix(5, 0),
		time.Unix(6, 0),
		time.Unix(7, 0),
		time.Unix(8, 0),
		time.Unix(9, 0),
	}
	dataObj = &obj{
		ID:        42,
		Name:      "bobheadxi",
		CreatedAt: time.Date(1998, 3, 11, 12, 0, 0, 0, time.UTC),
	}
	dataObjs = objs{
		dataObj,
		dataObj,
		dataObj,
		dataObj,
		dataObj,
		dataObj,
		dataObj,
		dataObj,
		dataObj,
	}
)

func fakeFields() []zap.Field {
	return []zap.Field{
		zap.Int("int", dataInts[0]),
		zap.Ints("ints", dataInts),
		zap.String("string", "hello world"),
		zap.Strings("strings", dataStrings),
		zap.Time("time", dataTimes[0]),
		zap.Times("times", dataTimes),
		zap.Object("obj", dataObj),
		zap.Object("another_obj", dataObj),
		zap.Array("users", dataObjs),
		zap.Error(errExample),
	}
}

func fakeMessage() string { return "#%v: the quick brown fox jumped over the sleeping dog" }
