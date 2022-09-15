package zapfields

import (
	"go.uber.org/zap"

	"github.com/ben0x539/errfields"
)

func With(l *zap.Logger, err error) *zap.Logger {
	fields := []zap.Field{zap.Error(err)}
	errfields.ProvideLogFields(err, func(key string, value any) {
		fields = append(fields, zap.Any(key, value))
	})

	return l.With(fields...)
}
