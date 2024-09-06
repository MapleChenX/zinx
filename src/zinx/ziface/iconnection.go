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

	/**
	 * 发送消息给客户端
	 * @param id 消息ID/cmd
	 * @param data 消息内容
	 */
	SendData(id uint32, data []byte) error

	// set property
	SetProperty(key string, value interface{})
	// get property
	GetProperty(key string) (interface{}, error)
	// remove property
	RemoveProperty(key string)
}

// define a function to handle connection business
type HandleFunc func(*net.TCPConn, []byte, int) error
