package opc

import (
	"bytes"
	"net"
	"testing"
	"time"
)

// This struct is used to mock out network connections
// So that we can test connection network operations accordingly
type MockConn struct {
	payload []byte
}

// Read a single byte off of the payload into the passed in byte array.
func (m *MockConn) Read(b []byte) (n int, err error) {
	b[0] = m.payload[0]
	m.payload = m.payload[len(b):]
	return len(b), nil
}

// Mocked out implementations of net.Conn.
func (m *MockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (m *MockConn) Close() error                       { return nil }
func (m *MockConn) LocalAddr() net.Addr                { return nil }
func (m *MockConn) RemoteAddr() net.Addr               { return nil }
func (m *MockConn) SetDeadline(t time.Time) error      { return nil }
func (m *MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *MockConn) SetWriteDeadline(t time.Time) error { return nil }

// This struct is used to mock out a device implementation
// such that we can test server operations accordingly
type MockDevice struct {
	channel uint8
}

func (md *MockDevice) Write(m *Message) error {
	return nil
}

func (md *MockDevice) Channel() uint8 {
	return 0
}

func TestRegisterDevice(t *testing.T) {
	s := NewServer()
	d := &MockDevice{channel: 1}
	s.RegisterDevice(d)

	if _, ok := s.Devs[d.Channel()]; !ok {
		t.Errorf("Expected Device to be registered")
	}
}

func TestUnregisterDevice(t *testing.T) {
	s := NewServer()
	d := &MockDevice{channel: 1}
	s.RegisterDevice(d)
	s.UnregisterDevice(d)

	if _, ok := s.Devs[d.Channel()]; ok {
		t.Errorf("Expected Device to be unregistered after registering it")
	}
}

func TestReadOpc(t *testing.T) {
	payload := []byte{255, 0, 0, 3, 1, 2, 3}
	m := &MockConn{payload: payload}

	msg, err := ReadOpc(m)
	if err != nil {
		t.Errorf("Encountered an error when reading a valid Message")
	}

	if bytes.Compare(msg.ByteArray(), payload) != 0 {
		t.Errorf("Recieved a mismatched message when reading from a Mocked Connection")
	}
}
