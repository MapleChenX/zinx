package apis

import (
	"fmt"
	"go_code/src/mmo_game_server/core"
	"go_code/src/mmo_game_server/pb"
	"go_code/src/zinx/ziface"
	"go_code/src/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	// 从客户端发送的数据中解析出msg
	msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Unmarshal msg err:", err)
		return
	}

	// 获取发送消息的玩家ID
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid err:", err)
		return
	}

	// 根据pid得到player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 将消息广播给其他全部在线玩家
	player.Talk(msg.Content)
}
