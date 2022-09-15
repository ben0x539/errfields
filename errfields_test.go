package errfields

import (
	"errors"
	"fmt"
	"testing"
)

type kvRecorder map[string]any

func (k kvRecorder) CollectLogFields(err error) {
	ProvideLogFields(err, func(key string, value any) {
		k[key] = value
	})
}

func (k kvRecorder) Require(t *testing.T, key string, want any) {
	have, ok := k[key]
	if !ok {
		t.Errorf("value missing for key=%v", key)
		return
	}

	if have != want {
		t.Errorf("wrong value for key=%v\nwant=%v\nhave=%v", key, want, have)
	}
}

func TestBasic(t *testing.T) {
	var err error = &withField{key: "my-field", value: 42}
	fields := kvRecorder{}
	fields.CollectLogFields(err)
	fields.Require(t, "my-field", 42)
}

func TestWrap(t *testing.T) {
	var err error = &withField{key: "inner", value: 1}
	err = &withField{key: "outer", value: 2, error: err}
	fields := kvRecorder{}
	fields.CollectLogFields(err)
	fields.Require(t, "inner", 1)
	fields.Require(t, "outer", 2)
}

func TestWrapOther(t *testing.T) {
	err := fmt.Errorf("innermost")
	err = &withField{key: "inner", value: 1, error: err}
	err = fmt.Errorf("middle: %w", err)
	err = &withField{key: "outer", value: 2, error: err}
	err = fmt.Errorf("outermost: %w", err)
	fields := kvRecorder{}
	fields.CollectLogFields(err)
	fields.Require(t, "inner", 1)
	fields.Require(t, "outer", 2)
}

func TestWrapOrder(t *testing.T) {
	var err error = &withField{key: "my-field", value: 1}
	err = &withField{key: "my-field", value: 2, error: err}
	fields := kvRecorder{}
	fields.CollectLogFields(err)
	fields.Require(t, "my-field", 1)
}

func TestMessage(t *testing.T) {
	err := fmt.Errorf("innermost")
	err = &withField{key: "inner", value: 1, error: err}
	err = fmt.Errorf("middle: %w", err)
	err = &withField{key: "outer", value: 2, error: err}
	err = fmt.Errorf("outermost: %w", err)

	want := "outermost: middle: innermost"
	have := err.Error()
	if want != have {
		t.Errorf("wrong string repr\nwant=%v\nhave=%v", want, have)
	}
}

func TestRequestOrder(t *testing.T) {
	var err error = &withField{key: "my-field", value: 1}
	err = &withField{key: "my-field", value: 2, error: err}

	want := 1
	have := RequestLogField(err, "my-field")
	if want != have {
		t.Errorf("wrong returned value\nwant=%v\nhave=%v", want, have)
	}
}

func TestAddZero(t *testing.T) {
	if err := Add(nil, "my-field", 42); err != nil {
		t.Errorf("expected nil error to remain nil, got %v", err)
	}
}

func TestAddOne(t *testing.T) {
	err := Add(errors.New("innermost"), "my-field", 42)
	if err == nil {
		t.Errorf("expected non-nil error to remain nil, nil")
		return
	}

	fields := kvRecorder{}
	fields.CollectLogFields(err)
	fields.Require(t, "my-field", 42)
}

func TestAddTwo(t *testing.T) {
	err := Add(errors.New("innermost"), "inner", 1)
	err = Add(err, "outer", 2)
	if err == nil {
		t.Errorf("expected non-nil error to remain nil, nil")
		return
	}

	fields := kvRecorder{}
	fields.CollectLogFields(err)
	fields.Require(t, "inner", 1)
	fields.Require(t, "outer", 2)
}

func TestAddMany(t *testing.T) {
	err := errors.New("innermost")
	for i := 0; i < 5; i++ {
		err = Add(err, fmt.Sprintf("field-%d", i), "abcde"[i])
		if err == nil {
			t.Errorf("expected non-nil error to remain nil, nil")
			return
		}

		err = fmt.Errorf("level %v: %w", i, err)
		if err == nil {
			t.Errorf("expected non-nil error to remain nil, nil")
			return
		}
	}

	fields := kvRecorder{}
	fields.CollectLogFields(err)
	fields.Require(t, "field-0", byte('a'))
	fields.Require(t, "field-1", byte('b'))
	fields.Require(t, "field-2", byte('c'))
	fields.Require(t, "field-3", byte('d'))
	fields.Require(t, "field-4", byte('e'))

	want := "level 4: level 3: level 2: level 1: level 0: innermost"
	have := err.Error()
	if want != have {
		t.Errorf("wrong string repr\nwant=%v\nhave=%v", want, have)
	}
}
