package errfields

import (
	"errors"
)

type Provider interface {
	ProvideLogFields(f func(key string, value any))
}

func ProvideLogFields(err error, f func(key string, value any)) {
	for err != nil {
		if p, ok := err.(Provider); ok {
			p.ProvideLogFields(f)
		}

		inner := errors.Unwrap(err)
		if inner == err {
			break
		}

		err = inner
	}
}

func RequestLogField(err error, desiredKey string) any {
	var result any

	ProvideLogFields(err, func(key string, value any) {
		if key == desiredKey {
			result = value
		}
	})

	return result
}
