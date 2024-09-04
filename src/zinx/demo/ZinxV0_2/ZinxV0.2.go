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
	fmt.Println("[default router] recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	err := request.GetConnection().SendData(1, []byte("ping...ping...ping测试 默认\n"))
	if err != nil {
		return
	}
}

type Router2 struct {
	znet.BaseRouter
}

func (pr *Router2) Handle(request ziface.IRequest) {
	fmt.Println("[2 router] recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	err := request.GetConnection().SendData(1, []byte("ping...ping...ping测试 2\n"))
	if err != nil {
		return
	}
}

type Router3 struct {
	znet.BaseRouter
}

func (pr *Router3) Handle(request ziface.IRequest) {
	fmt.Println("[3 router] recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	err := request.GetConnection().SendData(1, []byte("ping...ping...ping测试 3\n"))
	if err != nil {
		return
	}
}

type Router4 struct {
	znet.BaseRouter
}

func (pr *Router4) Handle(request ziface.IRequest) {
	fmt.Println("[4 router] recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	err := request.GetConnection().SendData(1, []byte("ping...ping...ping测试 4\n"))
	if err != nil {
		return
	}
}

func main() {
	server := znet.NewServer()

	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &Router2{})
	server.AddRouter(2, &Router3{})
	server.AddRouter(3, &Router4{})

	server.Serve()
}
