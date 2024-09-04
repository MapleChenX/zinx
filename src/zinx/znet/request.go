package znet

import "go_code/src/zinx/ziface"

type Request struct {
	// connection
	conn ziface.IConnection

	// request data
	msg ziface.IMessage
}

// get current connection
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// get request data
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// get request message ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
