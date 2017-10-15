package opc

import (
	"bytes"
	"testing"
)

func TestNewMessage(t *testing.T) {
	m := NewMessage(0)

	if m.Channel != 0 {
		t.Errorf("Unexpected channel value after initialization.")
	}

	if m.Command != SetPixelColorsCmd {
		t.Errorf("Unexpected command value after initialization.")
	}

	if m.HighLen != 0 {
		t.Errorf("Unexpected high length byte value after initialization.")
	}

	if m.LowLen != 0 {
		t.Errorf("Unexpected low length byte value after initialization.")
	}
}

func TestSetPixelColor(t *testing.T) {
	m := NewMessage(0)
	m.SetPixelColor(0, uint8(255), uint8(254), uint8(253))

	if m.Data[0] != 255 {
		t.Errorf("Did not set pixel 0's Red value correctly")
	}

	if m.Data[1] != 254 {
		t.Errorf("Did not set pixel 0's Green value correctly")
	}

	if m.Data[2] != 253 {
		t.Errorf("Did not set pixel 0's Blue value correctly")
	}
}

func TestSetLength(t *testing.T) {
	m := NewMessage(0)
	m.SetLength(10)
	if uint64(m.LowLen) != uint64(10) {
		t.Errorf("Expected a call to SetLength() to set the Message length to 10.")
	}

	m.SetLength(uint16(MaxMessageSize))
	if m.Length() != uint16(MaxMessageSize) {
		t.Errorf("Expected setting length to MAX_MESSAGE_SIZE to be reflected correctly.")
	}
}

func TestLength(t *testing.T) {
	m := NewMessage(0)
	m.LowLen = byte(10)
	v := m.Length()
	if v != uint16(10) {
		t.Errorf("Expected a call to Length() to return 10 after manually setting it to 10.")
	}
}

func TestIsValid(t *testing.T) {
	m := NewMessage(0)
	if m.IsValid() {
		t.Errorf("A Message should not be valid after initializing it.")
	}

	m.Data = make([]byte, 9)
	m.SetLength(9)

	if !m.IsValid() {
		t.Errorf("A Message should not be invalid after inserting data into its byte array and explicitly setting its length")
	}
}

func TestIsBroadcast(t *testing.T) {
	m := NewMessage(255)
	if m.IsBroadcast() {
		t.Errorf("A Message with a channel set to 255 should not be a Broadcast.")
	}

	m.Channel = byte(0)
	if !m.IsBroadcast() {
		t.Errorf("A Message with a channel set to 0 should be a Broadcast.")
	}
}

func TestByteArray(t *testing.T) {
	m := NewMessage(255)
	m.SetPixelColor(0, 1, 2, 3)
	m.SetLength(3)

	data := m.ByteArray()
	if bytes.Compare(data, []byte{255, 0, 0, 3, 1, 2, 3}) != 0 {
		t.Errorf("Unexpected message after converting to ByteArray got: %v", data)
	}
}
