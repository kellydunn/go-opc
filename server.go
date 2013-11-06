package opc

import (
	_ "fmt"
	"net"
)

// This struct describes an OPC server,
// which keeps track of all connected OPC devices
// as well as a channel of incoming messages from all connected clients
type Server struct {
	devs     map[uint8]Device
	messages chan *Message
}

// Creates and returns a new opc.Server.
// Accepts a list of usb product IDs in which to send opc messages to.
func NewServer() *Server {
	return &Server{devs: make(map[uint8]Device), messages: make(chan *Message)}
}

// Registers the passed in device to the OPC server
func (s *Server) RegisterDevice(dev Device) {
	s.devs[dev.Channel()] = dev
}

// Unregisters the passed in device from the OPC server
func (s *Server) UnregisterDevice(dev Device) {
	delete(s.devs, dev.Channel())
}

// Listens on the passed in port with the passed in protocol,
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

		go s.handleConn(conn)
	}
}

// Reads off OPC messages from the passed in connection
// until the connection breaks.
// Appends all valid messages onto the message channel
func (s *Server) handleConn(conn net.Conn) {
	for {
		msg, err := s.readOpc(conn)
		if err != nil {
			// If we encounter an error reading from the connection,
			// "break" out of the loop and stop reading.
			//
			// TODO find some way of maybe alerting to the client
			//      that an error occured
			break
		}

		s.messages <- msg
	}
}

// Reads and returns a single OPC message from the passed in connection.
func (s *Server) readOpc(conn net.Conn) (*Message, error) {
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
			m.channel = buf[0]
		case 2:
			m.command = buf[0]
		case 3:
			m.highLen = buf[0]
		case 4:
			m.lowLen = buf[0]
		default:
			m.data[bytesRead-5] = buf[0]

			if bytesRead-4 == m.Length() {
				m.data = m.data[:m.Length()]
			}
		}
	}

	return m, nil
}

// Dispatches the passed in message to all applicable devices.
// If the message is of a Broadcast type, it sends it to all connected devices
// Otherwise, it sends it to the specified device.
func (s *Server) dispatch(m *Message) {
	if m.IsBroadcast() {
		// Broadcast the message to all registered devices
		for i := range s.devs {
			s.devs[i].Write(m)
		}

	} else {
		// Otherwise write to the device specified by the message's channel
		//fmt.Printf("Attempting to write to device at channel:%d\n", m.channel)
		s.devs[m.channel].Write(m)
	}
}

// Processes all pending messages indefinitely
func (s *Server) Process() {
	for {
		msg := <-s.messages
		s.dispatch(msg)
	}
}
