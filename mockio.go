//	Package mockio provides a mock ReaderWriterCloser to test your software's i/o stream handling.
package mockio

/*
Package mockio implements a ReaderWriterCloser for mocking a ReaderWriterCloser
(such as a serial device).

You can program it to expect certain sequences of bytes and provide a response.

Alternatively, you can provide a function to determine if a sequence of bytes
written to the mock device matches your expectations, and send a response.

Or, you can create your own structure matching the Expect interface to handle
the incoming bytes and respond accordingly.
*/

import (
	"bytes"
	"time"
)

// MockIO provides a mock I/O ReadWriteCloser to test your software's i/o stream handling.
type MockIO struct {
	buffer  *bytes.Buffer
	holding []byte
	timer   chan time.Duration
	Expects []Expect
}

// NewMockIO constructs a new MockIO.
func NewMockIO() *MockIO {
	result := &MockIO{
		buffer:  bytes.NewBuffer([]byte{}),
		holding: []byte{},
		Expects: []Expect{},
		timer:   make(chan time.Duration, 5),
	}
	return result
}

// Read implements the Reader interface for the MockIO stream.
// Use this with your software to test that it reads correctly.
func (m *MockIO) Read(data []byte) (n int, err error) {
	d := <-m.timer
	<-time.NewTimer(d).C
	return m.buffer.Read(data)
}

// Write implements the Writer interface for the MockIO stream.
// Use this with your software to test that it writes correctly.
func (m *MockIO) Write(data []byte) (n int, err error) {
	m.holding = append(m.holding, data...)
	respond := []byte{}
	dur := time.Millisecond
	for _, test := range m.Expects {
		response, count, ok := test.Match(m.holding)
		if !ok {
			continue
		}
		dur = test.Duration()
		respond = append(respond, response...)
		m.holding = m.holding[count:]
		break
	}
	m.timer <- dur
	m.buffer.Write(respond)
	return len(data), nil
}

// Send writes data through the mock serial device to the Read function.  It is
// meant to mimic the serial device sending data not in response to a Write,
// but on its own (perhaps due to an event on the hardware connected to the
// serial device you're mocking).
func (m *MockIO) Send(data []byte, wait time.Duration) {
	m.timer <- wait
	m.buffer.Write(data)
}

// Expect adds a new Expect item to a list of things for the MockIO Write to expect.
func (m *MockIO) Expect(exp Expect) {
	m.Expects = append(m.Expects, exp)
}

// ClearExpectations removes all items stored in MockIO's Expect buffer.
func (m *MockIO) ClearExpectations() {
	m.Expects = []Expect{}
}

// Close implements the Closer interface for the MockIO stream.
func (m *MockIO) Close() (err error) {
	m.buffer.Reset()
	m.holding = []byte{}
	return err
}
