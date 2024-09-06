package main

import (
	"go_code/src/mmo_game_server/core"
	"go_code/src/zinx/ziface"
	"go_code/src/zinx/znet"
)

func DoConnectionBegin(conn ziface.IConnection) {
	// 给客户端发送id为1的消息
	player := core.NewPlayer(conn)
	player.SyncPid()

	// 给客户端发送id为200的消息
	player.BroadCastStartPosition()
}

func DoConnectionEnd(conn ziface.IConnection) {

}

func main() {
	// 1 创建一个server句柄
	server := znet.NewServer()

	// 2 注册连接Hook钩子函数
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionEnd)

	// 3 添加router
	server.AddRouter(1, &znet.BaseRouter{})

	// 4 启动server
	server.Serve()
}
