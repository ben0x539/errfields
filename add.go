package errfields

type withField struct {
	error

	key   string
	value any
}

type withFields struct {
	error

	fields []kv
}

type kv struct {
	key   string
	value any
}

func (w *withField) Unwrap() error { return w.error }

func (w *withField) ProvideLogFields(f func(key string, value any)) {
	f(w.key, w.value)
}

func (w *withFields) Unwrap() error { return w.error }

func (w *withFields) ProvideLogFields(f func(key string, value any)) {
	for i := len(w.fields) - 1; i >= 0; i-- {
		field := w.fields[i]
		f(field.key, field.value)
	}
}

func Add(err error, key string, value any) error {
	if w, ok := err.(withFields); ok {
		w.fields = append(w.fields, kv{key, value})
		return w
	}

	if w, ok := err.(withField); ok {
		return &withFields{
			error: w.error,
			fields: []kv{
				{w.key, w.value},
				{key, value},
			},
		}
	}

	if err != nil {
		return &withField{err, key, value}
	}

	return nil
}
