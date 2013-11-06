package opc

const (
	DEFAULT_OPC_PORT = "7890"
)

func ListenAndServe() {
	s := NewServer()
	go s.ListenOnPort("tcp", DEFAULT_OPC_PORT)
	go s.Process()
	select {}
}
