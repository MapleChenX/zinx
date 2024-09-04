package main

import (
	"fmt"
	"go_code/src/zinx/ziface"
	"go_code/src/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	err := request.GetConnection().SendData(1, []byte("ping...ping...ping测试\n"))
	if err != nil {
		return
	}
}

func main() {
	server := znet.NewServer()

	server.AddRouter(&PingRouter{})

	server.Serve()
}
