package apis

import (
	"go_code/src/mmo_game_server/core"
	"go_code/src/mmo_game_server/pb"
	"go_code/src/zinx/ziface"
	"go_code/src/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type MoveApi struct {
	znet.BaseRouter
}

func (ma *MoveApi) Handle(request ziface.IRequest) {
	position := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), position)
	if err != nil {
		return
	}

	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		return
	}

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.UpdatePosition(position.X, position.Y, position.Z, position.V)
}
