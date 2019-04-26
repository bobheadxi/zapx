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

type fieldSet struct {
	fields []zap.Field
}

func (f fieldSet) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, field := range f.fields {
		field.AddTo(enc)
	}
	return nil
}

func (f fieldSet) Fields() []zap.Field { return f.fields }

// Fields allows logging of sets of fields as a log object
func Fields(key string, fields ...zap.Field) zap.Field {
	return zap.Field{
		Key:       key,
		Type:      zapcore.ObjectMarshalerType,
		Interface: fieldSet{fields: fields},
	}
}
