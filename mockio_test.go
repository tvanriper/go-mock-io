package mockio

import (
	"io"
	"log"
	"testing"
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
	m.Expect(NewExpectBytes([]byte{1, 2}, exp))
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
	m.Expect(NewExpectFunc(test, exp))
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
	m.Expect(NewExpectBytes([]byte{1, 2}, exp))
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
	m.Expect(NewExpectBytes([]byte{1, 2}, exp))
	l := len(m.Expects)
	if l != 1 {
		t.Errorf("expected 1 item, but found %d", l)
	}
	m.ClearExpectations()
	l = len(m.Expects)
	if l != 0 {
		t.Errorf("expected 0 items, but found %d", l)
	}

	m.Expect(NewExpectBytes([]byte{1, 2}, exp))
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
