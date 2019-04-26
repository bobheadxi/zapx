package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// FieldSetMarshaller is the underlying type used by zapx.Fields()
type FieldSetMarshaller interface {
	zapcore.ObjectMarshaler
	Fields() []zap.Field
}

type fieldSet []zap.Field

func (f fieldSet) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for i := range f {
		f[i].AddTo(enc)
	}
	return nil
}

func (f fieldSet) Fields() []zap.Field { return f }

// FieldSet allows logging of sets of fields as a log object
// TODO: not as performant as testdata.Obj for some reason, need to look into why
func FieldSet(key string, fields ...zap.Field) zap.Field {
	return zap.Field{
		Key:       key,
		Type:      zapcore.ObjectMarshalerType,
		Interface: fieldSet(fields),
	}
}
