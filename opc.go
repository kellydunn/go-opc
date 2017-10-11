package opc

const (
	// DefaultOpcPort is the default Open Pixel Control protocol port.
	DefaultOpcPort = "7890"
)

// ListenAndServe creates a new OPC server and listens indefinently.
func ListenAndServe() {
	s := NewServer()
	go s.ListenOnPort("tcp", DefaultOpcPort)
	go s.Process()
	select {}
}
