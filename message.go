package opc

const (
	SET_PIXEL_COLORS  = 0x00
	SYSTEM_EXCLUSIVE  = 0xFF
	HEADER_BYTES      = 4
	BROADCAST_CHANNEL = 0
	MAX_MESSAGE_SIZE  = 0xFFFF
)

// This struct describes a single message
// that follows the OPC protocol
type Message struct {
	channel byte
	command byte
	highLen byte
	lowLen  byte
	data    []byte
}

// Creates and returns a pointer to a new message that is to be sent
// to the passed in channel
func NewMessage(channel uint8) *Message {
	return &Message{channel: channel, data: make([]byte, MAX_MESSAGE_SIZE)}
}

// Sets the pixel color of the passed in pixel
// to the passed in red, green, and blue colors, respectively for this message
func (m *Message) SetPixelColor(pixel int, r uint8, g uint8, b uint8) {
	index := (3 * pixel)
	m.data[index] = r
	m.data[index+1] = g
	m.data[index+2] = b
}

// Specifies that this message is a System Exclusive Message
// and populates data accordingly
func (m *Message) SystemExclusive(systemId []byte, data []byte) {
	m.command = SYSTEM_EXCLUSIVE
	m.data = systemId
	for i := 0; i < len(data); i++ {
		m.data = append(m.data, data[i])
	}
}

// Sets the length of this message by splitting
// the passed in length into high and low length bytes.
func (m *Message) SetLength(length uint16) {
	m.highLen = byte(length >> 8)
	m.lowLen = byte(length)
}

// Returns the length of the message.
// The length of the message is respresented by combining
// the high and low length bytes of this message.
func (m *Message) Length() uint16 {
	return (uint16(m.highLen) << 8) | uint16(m.lowLen)
}

// Returns whether or not this message is valid or not.
// Validity is determined as whether or not the Length of the message
// corresponds with the number of data bytes in the message
func (m *Message) IsValid() bool {
	return m.Length() == uint16(len(m.data))
}

// Returns whether or not this message is a Broadcast message.
func (m *Message) IsBroadcast() bool {
	return m.channel == BROADCAST_CHANNEL
}

// Returns a byte array representation of this message.
func (m *Message) ByteArray() []byte {
	data := []byte{}
	data = append(data, m.channel, m.command, m.highLen, m.lowLen)
	for i := uint16(0); i < m.Length(); i++ {
		data = append(data, m.data[i])
	}
	return data
}
