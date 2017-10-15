package opc

// OPC constants.
const (
	SetPixelColorsCmd  = 0x00   // SetPixelColorsCmd is the command for setting pixelcolor.
	SystemExclusiveCmd = 0xFF   // SystemExclusiveCmd is denotes a command for a specific system.
	HeaderBytes        = 4      // HeaderBytes is the number of bytes an OPC header consists of.
	BroadcastChannel   = 0      // BroadcastChannel is the broadcast channel number.
	MaxMessageSize     = 0xFFFF // MaxMessageSize is the max message size for the OPC protocol.
)

// Message describes a single message
// that follows the OPC protocol
type Message struct {
	Channel byte
	Command byte
	HighLen byte
	LowLen  byte
	Data    []byte
}

// NewMessage creates and returns a pointer to a new message that is to be sent
// to the passed in channel
func NewMessage(channel uint8) *Message {
	return &Message{Channel: channel, Data: make([]byte, MaxMessageSize)}
}

// SetPixelColor sets the pixel color of the passed in pixel
// to the passed in red, green, and blue colors, respectively for this message
func (m *Message) SetPixelColor(pixel int, r uint8, g uint8, b uint8) {
	index := (3 * pixel)
	m.Data[index] = r
	m.Data[index+1] = g
	m.Data[index+2] = b
}

// SystemExclusive specifies that this message is a System Exclusive Message
// and populates data accordingly
func (m *Message) SystemExclusive(systemID []byte, data []byte) {
	m.Command = SystemExclusiveCmd
	m.Data = systemID
	for i := 0; i < len(data); i++ {
		m.Data = append(m.Data, data[i])
	}
}

// SetLength sets the length of this message by splitting
// the passed in length into high and low length bytes.
func (m *Message) SetLength(length uint16) {
	m.HighLen = byte(length >> 8)
	m.LowLen = byte(length)
}

// Length returns the length of the message.
// The length of the message is respresented by combining
// the high and low length bytes of this message.
func (m *Message) Length() uint16 {
	return (uint16(m.HighLen) << 8) | uint16(m.LowLen)
}

// Returns whether or not this message is valid or not.
// Validity is determined as whether or not the Length of the message
// corresponds with the number of data bytes in the message
func (m *Message) IsValid() bool {
	return m.Length() == uint16(len(m.Data))
}

// Returns whether or not this message is a Broadcast message.
func (m *Message) IsBroadcast() bool {
	return m.Channel == BroadcastChannel
}

// ByteArray returns a byte array representation of this message.
func (m *Message) ByteArray() []byte {
	data := []byte{}
	data = append(data, m.Channel, m.Command, m.HighLen, m.LowLen)
	for i := uint16(0); i < m.Length(); i++ {
		data = append(data, m.Data[i])
	}
	return data
}
