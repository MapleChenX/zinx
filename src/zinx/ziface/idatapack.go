package ziface

import "net"

type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
	data2msg(conn net.Conn) (IMessage, error)
}
