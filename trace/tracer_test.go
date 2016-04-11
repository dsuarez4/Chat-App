package trace

import (
	"testing"
	"bytes"
)

func TestNew(t *testing.T) {

	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Return from 'new' should not be nil")
	} else {
		tracer.Trace("Hello Trace package.")
		if buf.String() != "Hello Trace package.\n" {
			t.Errorf("Trace should not write '%s'.", buf.String())
		}
	}
}

