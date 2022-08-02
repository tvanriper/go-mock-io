package mockio

import "time"

// Expect provides an interface for the mock io stream to use when matching incoming
// information.
type Expect interface {
	Match(b []byte) (response []byte, count int, ok bool) // Match provides a response to a given set of bytes
	Duration() time.Duration                              // Duraction provides the amount of time to wait before responding to a read request for this response
}

// ExpectBytes provides a kind of Expect that precisely matches a sequence of bytes to
// respond with another sequence of bytes.
type ExpectBytes struct {
	Expect       []byte
	Respond      []byte
	WaitDuration time.Duration
}

// Match implements the Expect interface for ExpectBytes.
func (e *ExpectBytes) Match(b []byte) (response []byte, count int, ok bool) {
	l := len(e.Expect)
	ok = true
	for i := 0; i < len(b); i++ {
		if i == l {
			break
		}
		if b[i] != e.Expect[i] {
			ok = false
			break
		}
	}

	if ok {
		response = append(response, e.Respond...)
		count = l
	}
	return response, count, ok
}

// Duration responds with how long to wait before responding to a read request.
func (e *ExpectBytes) Duration() time.Duration {
	return e.WaitDuration
}

// NewExpectBytes provides a convenience constructor for ExpectBytes.
func NewExpectBytes(expect []byte, respond []byte, wait time.Duration) *ExpectBytes {
	return &ExpectBytes{
		Expect:       expect,
		Respond:      respond,
		WaitDuration: wait,
	}
}

// ExpectFuncTest describes a test that matches a sequence of bytes written to the mock
// io stream.
// Always start the match from b[0].  Use count to indicate how many bytes in b were
// consumed in the test.  Set ok to true if there was a successful match.
type ExpectFuncTest func(b []byte) (count int, ok bool)

// ExpectFunc provides a kind of Expect that tests a sequence of bytes with a function.
type ExpectFunc struct {
	Test         ExpectFuncTest
	Respond      []byte
	WaitDuration time.Duration
}

// Match implements the Expect interface for ExpectFunc.
func (e *ExpectFunc) Match(b []byte) (response []byte, count int, ok bool) {
	count, ok = e.Test(b)
	if ok {
		response = append(response, e.Respond...)
	}
	return response, count, ok
}

// Duration provides the amount of time to wait before responding to a read request.
func (e *ExpectFunc) Duration() time.Duration {
	return e.WaitDuration
}

// NewExpectFunc provides a convenience constructor for ExpectFunc.
func NewExpectFunc(fn ExpectFuncTest, response []byte, wait time.Duration) *ExpectFunc {
	return &ExpectFunc{
		Test:         fn,
		Respond:      response,
		WaitDuration: wait,
	}
}
