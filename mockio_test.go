package mockio

import (
	"io"
	"log"
	"math/rand"
	"testing"
	"time"
)

func MatchBytes(l []byte, r []byte) bool {
	ll := len(l)
	lr := len(r)
	if ll != lr {
		return false
	}
	for i := 0; i < ll; i++ {
		if l[i] != r[i] {
			return false
		}
	}
	return true
}

func TestMockIO(t *testing.T) {
	m := NewMockIO()
	exp := []byte{3, 4}
	m.Expect(NewExpectBytes([]byte{1, 2}, exp, 0))
	m.Write([]byte{1, 2})
	b := make([]byte, 2)
	n, err := m.Read(b)
	if err != nil {
		t.Error(err)
	}
	if n != 2 {
		t.Errorf("expected %d, received %d", len(exp), n)
	}

	if !MatchBytes(b, exp) {
		t.Errorf("expected %#v but received %#v", exp, b)
	}

	m.Write([]byte{1, 2, 1, 2})

	b = make([]byte, 2)
	n, err = m.Read(b)
	if err != nil {
		t.Error(err)
	}
	if n != 2 {
		t.Errorf("expected %d, received %d", len(exp), n)
	}

	if !MatchBytes(b, exp) {
		t.Errorf("expected %#v but received %#v", exp, b)
	}

	l := len(m.holding)
	if l != 2 {
		t.Errorf("expected 2 bytes holding, but found %d", l)
	}
}

func TestMockIOFunc(t *testing.T) {
	m := NewMockIO()
	test := func(b []byte) (count int, ok bool) {
		if len(b) != 2 {
			return 0, false
		}
		if b[0] != 0 && b[1] != 1 {
			return 0, false
		}
		return 2, true
	}
	exp := []byte{3, 4}
	m.Expect(NewExpectFunc(test, exp, 0))
	m.Write([]byte{0, 1})
	b := make([]byte, 2)
	n, err := m.Read(b)
	if err != nil {
		t.Error(err)
	}
	if n != 2 {
		t.Errorf("expected %d, received %d", len(exp), n)
	}

	if !MatchBytes(b, exp) {
		t.Errorf("expected %#v but received %#v", exp, b)
	}

	log.Printf("received: %#v\n", b)
}

func TestMockIOFail(t *testing.T) {
	m := NewMockIO()
	exp := []byte{3, 4}
	m.Expect(NewExpectBytes([]byte{1, 2}, exp, 0))
	m.Write([]byte{0, 2})
	b := make([]byte, 2)
	n, err := m.Read(b)
	if err == nil {
		t.Errorf("Expected failure")
	}
	if err != io.EOF {
		t.Errorf("expected EOF but got: %s", err)
	}
	if n != 0 {
		t.Errorf("expected no data, received %d", n)
	}
	l := len(m.buffer.Bytes())
	if l != 0 {
		t.Errorf("expected no bytes in buffer, but found %d: %#v", l, m.buffer.Bytes())
	}

	m.Write([]byte{1, 4})
	b = make([]byte, 2)
	n, err = m.Read(b)
	if err == nil {
		t.Errorf("Expected failure")
	}
	if err != io.EOF {
		t.Errorf("expected EOF but got: %s", err)
	}
	if n != 0 {
		t.Errorf("expected no data, received %d", n)
	}
	l = len(m.buffer.Bytes())
	if l != 0 {
		t.Errorf("expected no bytes in buffer, but found %d: %#v", l, m.buffer.Bytes())
	}
}

func TestMockClear(t *testing.T) {
	m := NewMockIO()
	exp := []byte{3, 4}
	m.Expect(NewExpectBytes([]byte{1, 2}, exp, 0))
	l := len(m.Expects)
	if l != 1 {
		t.Errorf("expected 1 item, but found %d", l)
	}
	m.ClearExpectations()
	l = len(m.Expects)
	if l != 0 {
		t.Errorf("expected 0 items, but found %d", l)
	}

	m.Expect(NewExpectBytes([]byte{1, 2}, exp, 0))
	m.Write([]byte{3, 4})
	l = len(m.holding)
	if l != 2 {
		t.Errorf("holding should have 2 items, but found %d", l)
	}
	m.Close()
	l = len(m.holding)
	if l != 0 {
		t.Errorf("holding should have no items, but found %d", l)
	}
}

func TestMockDurations(t *testing.T) {
	for i := 0; i < 10; i++ {
		m := NewMockIO()
		test := func(b []byte) (count int, ok bool) {
			if len(b) != 2 {
				return 0, false
			}
			if b[0] != 0 && b[1] != 1 {
				return 0, false
			}
			return 2, true
		}
		exp := []byte{3, 4}
		e := NewExpectFunc(test, exp, time.Duration(rand.Intn(200))*time.Millisecond)
		m.Expect(e)
		start := time.Now()
		m.Write([]byte{0, 1})
		b := make([]byte, 2)
		n, err := m.Read(b)
		if err != nil {
			t.Error(err)
		}
		duration := time.Since(start)
		if n != 2 {
			t.Errorf("expected %d, received %d", len(exp), n)
		}

		if !MatchBytes(b, exp) {
			t.Errorf("expected %#v but received %#v", exp, b)
		}

		if duration < e.WaitDuration {
			t.Errorf("expected actual execution time to take longer than %s", e.WaitDuration.String())
		}

		log.Printf("received: %#v\n", b)
	}
}

func TestMockSend(t *testing.T) {
	m := NewMockIO()
	m.Send([]byte("Hi"), 0)
	b := make([]byte, 2)
	n, err := m.Read(b)
	if err != nil {
		t.Error(err)
	}
	if n != 2 {
		t.Errorf("expected 2, received %d", n)
	}
	if string(b) != "Hi" {
		t.Errorf("expected 'Hi' but received %#v", b)
	}

	// Testing 'duration'.
	b = make([]byte, 10)
	msg := "Hullo"
	start := time.Now()
	m.Send([]byte(msg), 100*time.Millisecond)
	n, err = m.Read(b)
	duration := time.Since(start)
	answer := string(b[:n])
	if err != nil {
		t.Error(err)
	} else {
		if duration < 100*time.Millisecond {
			t.Errorf("expected execution time to take longer than 100 milliseconds, but it took %s", duration.String())
		}
		if answer != msg {
			t.Errorf("expected %s but received %s", msg, answer)
		}
	}
}
