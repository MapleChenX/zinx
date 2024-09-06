package mmo_game_server

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

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("Do something when conn starts")
	conn.SendData(123, []byte("Do something when conn starts\n"))

	// 设置连接属性 --- 扩展内容测试
	conn.SetProperty("Name", "zinx")
	fmt.Println("Set Name = zinx")
}

func DoConnectionEnd(conn ziface.IConnection) {
	// 用户连接关闭日志
	fmt.Printf("connection %d/%s closed\n", conn.GetConnID(), conn.RemoteAddr())

	// 获取连接属性 --- 扩展内容测试
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
}

func main() {
	// 1 创建一个server句柄
	server := znet.NewServer()

	// 2 注册连接Hook钩子函数
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionEnd)

	// 3 添加router
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &Router2{})
	server.AddRouter(2, &Router3{})
	server.AddRouter(3, &Router4{})

	// 4 启动server
	server.Serve()
}
