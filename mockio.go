/*
	Package mockio provides a mock I/O ReadWriteCloser stream for testing your software's i/o stream handling.
	Example usage:

		func main() {
			s := NewMockIO()
			s.Expect(NewExpectBytes([]byte{1,2},[]byte{3,4}))
			s.Write([]byte{1,2})
			b := make([]byte,2)
			n, err := s.Read(b)
			if err != nil {
				panic(err)
			}
			fmt.Printf("n: %d, b: %#v", n, b)
		}

	which should print:

		n: 2, b: []byte{0x3, 0x4}

	You can use this to test things like serial I/O, or perhaps console I/O, without having
	to open a serial port or a console, allowing for automated testing in a controlled fashion.
*/
package mockio

import (
	"bytes"
	"time"
)

// MockIO provides a mock I/O ReadWriteCloser to test your software's i/o stream handling.
type MockIO struct {
	buffer  *bytes.Buffer
	holding []byte
	timer   <-chan time.Time
	Expects []Expect
}

// NewMockIO constructs a new MockIO.
func NewMockIO() *MockIO {
	return &MockIO{
		buffer:  bytes.NewBuffer([]byte{}),
		holding: []byte{},
		Expects: []Expect{},
	}
}

// Read implements the Reader interface for the MockIO stream.
// Use this with your software to test that it reads correctly.
func (m *MockIO) Read(data []byte) (n int, err error) {
	<-m.timer
	return m.buffer.Read(data)
}

// Write implements the Writer interface for the MockIO stream.
// Use this with your software to test that it writes correctly.
func (m *MockIO) Write(data []byte) (n int, err error) {
	m.holding = append(m.holding, data...)
	respond := []byte{}
	for _, test := range m.Expects {
		response, count, ok := test.Match(m.holding)
		if !ok {
			continue
		}
		m.timer = time.NewTimer(test.Duration()).C
		respond = append(respond, response...)
		m.holding = m.holding[count:]
		break
	}
	m.buffer.Write(respond)
	return len(data), nil
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
