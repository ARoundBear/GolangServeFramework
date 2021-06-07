package ziface

import "net"

//define the connection link module abstract layer
type IConnection interface {
	StartReader()
	//start link ready to work
	Start()

	//stop link end the linker work
	Stop()

	//get the link bind socket conn
	GetTCPConnection() *net.TCPConn

	//get Id of the link module
	GetConnID() uint32

	//get remote client ip port status pf tcp
	RemoteAddr() net.Addr

	// send data to remote client
	SendMsg(uint32, []byte) error
}

type HandleFunc func(*net.TCPConn, []byte, int) error
