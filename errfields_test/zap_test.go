package errfields_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ben0x539/errfields/zapfields"
)

func require(t *testing.T, buf []byte, key string, want any) {
	var val map[string]any
	err := json.Unmarshal(buf, &val)
	if err != nil {
		t.Errorf("couldn't unmarshal json log %q: %v", string(buf), err)
		return
	}

	have, ok := val[key]
	if !ok {
		t.Errorf("value missing for key=%v", key)
		return
	}

	if have != want {
		t.Errorf("wrong value for key=%v\nwant=%v\nhave=%v", key, want, have)
	}
}

func TestZap(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(enc, zapcore.AddSync(buf), zapcore.DebugLevel)
	err := &LookupError{lookupKey: "my-key"}
	log := zap.New(core)
	zapfields.WithErrorFields(log, err).Info("whatever")
	require(t, buf.Bytes(), "app-operation", "lookup")
	require(t, buf.Bytes(), "app-lookup-key", "my-key")
}
