package main

import (
	"fmt"
	"go_code/src/zinx/znet"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 封装消息
		pack := znet.NewDataPack()
		id := 1
		data, err := pack.Pack(znet.NewMessage(uint32(id), []byte("ZinxV0.2 client test message")))
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write conn err:", err)
			return
		}

		// 读取服务器返回的消息
		msg, err := pack.GetMsgFromConn(conn)
		if err != nil {
			fmt.Println("client unpack err:", err)
			return
		}

		fmt.Println("==> Recv MsgID: ", msg.GetMsgId(),
			", dataLen: ", msg.GetDataLen(),
			", data: ", string(msg.GetData()))

		time.Sleep(1 * time.Second)
		id++
	}
}
