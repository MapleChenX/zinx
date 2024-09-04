package ziface

/*
   把链接和请求的数据封装到一个 Request 中，作为一次请求的数据包
*/

type IRequest interface {
	// get current connection
	GetConnection() IConnection

	// get request data
	GetData() []byte

	// get request message ID
	GetMsgID() uint32
}
