package errfields

type LookupError struct {
	error

	lookupKey string
}

func (l *LookupError) Unwrap() error { return l.error }

func (l *LookupError) ProvideLogFields(f func(key string, value any)) {
	f("app-operation", "lookup")
	f("app-lookup-key", l.lookupKey)
}
