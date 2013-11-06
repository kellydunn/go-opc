package opc

// This interface describes the behavior of an OPC device.
// OPC devices should be able to write OPC messages to themselves
// as well as be able to announce a Channel in which they are listening on.
type Device interface {
	Write(*Message) error
	Channel() uint8
}
