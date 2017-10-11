package opc

import (
	"net"
)

// Server describes an OPC server,
// which keeps track of all connected OPC devices
// as well as a channel of incoming messages from all connected clients
type Server struct {
	Devs     map[uint8]Device
	Messages chan *Message
}

// NewServer creates and returns a new opc.Server.
// Accepts a list of usb product IDs in which to send opc messages to.
func NewServer() *Server {
	return &Server{Devs: make(map[uint8]Device), Messages: make(chan *Message)}
}

// RegisterDevice registers the passed in device to the OPC server
func (s *Server) RegisterDevice(dev Device) {
	s.Devs[dev.Channel()] = dev
}

// UnregisterDevice unregisters the passed in device from the OPC server
func (s *Server) UnregisterDevice(dev Device) {
	delete(s.Devs, dev.Channel())
}

// ListenOnPort listens on the passed in port with the passed in protocol,
// which in turn accepts incoming connections and handles them
// by issuing individual goroutines.
func (s *Server) ListenOnPort(protocol string, port string) {
	listener, listenerErr := net.Listen(protocol, port)
	if listenerErr != nil {
		panic(listenerErr)
	}

	for {
		conn, connErr := listener.Accept()
		if connErr != nil {
			panic(connErr)
		}

		go s.HandleConn(conn)
	}
}

// HandleConn reads off OPC messages from the passed in connection
// until the connection breaks.
// Appends all valid messages onto the message channel.
// ListenOnPort will accept clients and pass the processing on to
// this function.
func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()
	for {
		msg, err := ReadOpc(conn)
		if err != nil {
			// If we encounter an error reading from the connection,
			// "break" out of the loop and stop reading.
			break
		}

		s.Messages <- msg
	}
}

// ReadOpc reads and returns a single OPC message from the passed in connection.
func ReadOpc(conn net.Conn) (*Message, error) {
	buf := make([]byte, 1)
	bytesRead := uint16(0)
	m := NewMessage(0)

	for !m.IsValid() {
		_, err := conn.Read(buf)

		// Encountered an error in reading from connection!
		// Bail out with error message
		if err != nil {
			return nil, err
		}

		bytesRead++

		// Ignore first 4 bytes to account for HEADER_BYTES
		switch bytesRead {
		case 1:
			m.Channel = buf[0]
		case 2:
			m.Command = buf[0]
		case 3:
			m.HighLen = buf[0]
		case 4:
			m.LowLen = buf[0]
		default:
			m.Data[bytesRead-5] = buf[0]

			if bytesRead-4 == m.Length() {
				m.Data = m.Data[:m.Length()]
			}
		}
	}

	return m, nil
}

// Dispatch dispatches the passed in message to all applicable devices.
// If the message is of a Broadcast type, it sends it to all connected devices
// Otherwise, it sends it to the specified device.
// You can use Process to let the server automatically call this Dispatch for
// each incoming message. Or listen to the Messages channel on Server yourself
// and Dispatch it yourself.
func (s *Server) Dispatch(m *Message) {
	if m.IsBroadcast() {
		// Broadcast the message to all registered devices
		for i := range s.Devs {
			s.Devs[i].Write(m)
		}

	} else {
		// Otherwise write to the device specified by the message's channel
		s.Devs[m.Channel].Write(m)
	}
}

// Process processes all pending messages indefinitely.
func (s *Server) Process() {
	for {
		msg := <-s.Messages
		s.Dispatch(msg)
	}
}
