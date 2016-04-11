package trace

import "io"

// capital letters indicate user facing functions
type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) {}
