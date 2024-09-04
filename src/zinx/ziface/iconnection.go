package ziface

import "net"

type IConnection interface {
	// start connection
	Start()
	// stop connection
	Stop()
	// get connection ID
	GetConnID() uint32
	// get connection
	GetTCPConnection() *net.TCPConn
	// get remote client address
	RemoteAddr() net.Addr
	// send data to client
	Send(data []byte) error

	// send message to client
	SendMsg(msg IMessage) error
	SendData(id uint32, data []byte) error
}

// define a function to handle connection business
type HandleFunc func(*net.TCPConn, []byte, int) error
