module github.com/ben0x539/errfields/example

go 1.18

require (
	github.com/ben0x539/errfields v0.0.0-00010101000000-000000000000
	github.com/ben0x539/errfields/zapfields v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.23.0
)

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
)

replace github.com/ben0x539/errfields => ../

replace github.com/ben0x539/errfields/zapfields => ../zapfields
