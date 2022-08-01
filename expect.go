package mockio

import "time"

// Expect provides an interface for the mock io stream to use when matching incoming
// information.
type Expect interface {
	Match(b []byte) (response []byte, count int, ok bool)
	Duration() time.Duration
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

// Duration determines how long to wait to answer what is written.
func (e *ExpectBytes) Duration() time.Duration {
	return e.WaitDuration
}

// NewExpectBytes provides a convenience constructor for ExpectBytes.
func NewExpectBytes(expect []byte, respond []byte) *ExpectBytes {
	return &ExpectBytes{
		Expect:       expect,
		Respond:      respond,
		WaitDuration: 0,
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

func (e *ExpectFunc) Duration() time.Duration {
	return e.WaitDuration
}

// NewExpectFunc provides a convenience constructor for ExpectFunc.
func NewExpectFunc(fn ExpectFuncTest, response []byte) *ExpectFunc {
	return &ExpectFunc{
		Test:         fn,
		Respond:      response,
		WaitDuration: 0,
	}
}
