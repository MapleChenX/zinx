package main

import (
	"go_code/src/mmo_game_server/apis"
	"go_code/src/mmo_game_server/core"
	"go_code/src/zinx/ziface"
	"go_code/src/zinx/znet"
)

func DoConnectionBegin(conn ziface.IConnection) {
	player := core.NewPlayer(conn)

	// 给客户端发送id为1的消息
	player.SyncPid()

	// 给客户端发送id为200的消息
	player.BroadCastStartPosition()

	// 将新上线的玩家添加到世界管理器中
	core.WorldMgrObj.AddPlayer(player)

	// 将该连接绑定一个玩家ID
	conn.SetProperty("pid", player.PID)
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
	server.AddRouter(2, &apis.WorldChatApi{})

	// 4 启动server
	server.Serve()
}
