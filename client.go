package opc

import (
	"net"
)

// This struct represents an OPC client
// which is used to send OPC messages to an OPC server.
type Client struct {
	conn net.Conn
}

// Creates and returns a new Client
func NewClient() *Client {
	return &Client{}
}

// Connects the client to a server specified by
// the protocol string of either 'tcp' or 'udp', and the host location,
// which is a single string in the `url:port` format.
func (c *Client) Connect(protocol string, host string) error {
	conn, err := net.Dial(protocol, host)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

// Sends an OPC message from the Client to the Server connection.
func (c *Client) Send(m *Message) error {
	_, err := c.conn.Write(m.ByteArray())
	if err != nil {
		return err
	}

	return nil
}
