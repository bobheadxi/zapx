package testdata

import (
	"errors"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Obj is an object
type Obj struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

// MarshalLogObject implements zapcore.ObjectMarshaller
func (o *Obj) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", o.ID)
	enc.AddString("name", o.Name)
	enc.AddTime("created_at", o.CreatedAt)
	return nil
}

// Objs is an array of objects
type Objs []*Obj

// MarshalLogArray implements zapcore.ArrayMarshaller
func (o Objs) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	var err error
	for _, v := range o {
		err = multierr.Append(err, arr.AppendObject(v))
	}
	return err
}

var (
	// ErrExample is a sample err
	ErrExample = errors.New("hello world")
	// DataInts is sample data
	DataInts = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	// DataStrings is sample data
	DataStrings = []string{"the", "quick", "brown", "fox", "jumped", "over", "the", "sleeping", "dog", "amazing"}
	// DataTimes is sample data
	DataTimes = []time.Time{
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
	// DataObjs is sample data
	DataObjs = Objs{
		FakeObject(),
		FakeObject(),
		FakeObject(),
		FakeObject(),
		FakeObject(),
		FakeObject(),
		FakeObject(),
		FakeObject(),
		FakeObject(),
	}
)

// FakeFields produces a set of fake fields
func FakeFields() []zap.Field {
	return []zap.Field{
		zap.Int("int", DataInts[0]),
		zap.Ints("ints", DataInts),
		zap.String("string", "hello world"),
		zap.Strings("strings", DataStrings),
		zap.Time("time", DataTimes[0]),
		zap.Times("times", DataTimes),
		zap.Object("obj", FakeObject()),
		zap.Object("another_obj", FakeObject()),
		zap.Array("users", DataObjs),
		zap.Error(ErrExample),
	}
}

// FakeObject produces a fake object
func FakeObject() *Obj {
	return &Obj{
		ID:        42,
		Name:      "bobheadxi",
		CreatedAt: time.Date(1998, 3, 11, 12, 0, 0, 0, time.UTC),
	}
}

// FakeMessage produces a fake message
func FakeMessage() string { return "the quick brown fox jumped over the sleeping dog" }
