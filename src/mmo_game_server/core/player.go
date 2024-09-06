package core

import (
	"fmt"
	"go_code/src/mmo_game_server/pb"
	"go_code/src/zinx/ziface"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
)

type Player struct {
	PID  int32
	Conn ziface.IConnection

	// 坐标
	X float32
	Y float32
	Z float32
	V float32
}

var PidGen int32 = 100
var PidLock sync.RWMutex

// 一个玩家的对象
func NewPlayer(conn ziface.IConnection) *Player {
	PidLock.Lock()
	defer PidLock.Unlock()

	p := &Player{
		PID:  PidGen,
		Conn: conn,
		// 随机生成玩家的初始坐标
		// x和z才是地面坐标，y是高度
		X: float32(160 + rand.Intn(10)),
		Y: 0,
		Z: float32(160 + rand.Intn(10)),
		V: 0,
	}

	PidGen++
	return p
}

// 发送消息给客户端
func (p *Player) SendMsg(msgId uint32, msg proto.Message) {
	if p == nil || p.Conn == nil {
		fmt.Println("当前玩家为空！无法发送消息")
		return
	}

	// proto序列化数据
	data, err := proto.Marshal(msg)
	if err != nil {
		return
	}

	// tcp发送数据
	if err := p.Conn.SendData(msgId, data); err != nil {
		fmt.Println("Send data error: ", err)
		return
	}
}

// 同步玩家ID
func (p *Player) SyncPid() {
	// 组建MsgID:1消息
	msg := &pb.SyncPid{
		Pid: p.PID,
	}

	// 发送消息给客户端
	p.SendMsg(1, msg)
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	// 组建MsgID:200消息
	msg := &pb.BroadCast{
		Pid: p.PID,
		Tp:  2, // 2-广播位置信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				// 这里的坐标需要转换，因为客户端把X和Z坐标颠倒了
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMsg(200, msg)
}

// talk
func (p *Player) Talk(content string) {
	// 组建MsgID:200消息
	msg := &pb.BroadCast{
		Pid: p.PID,
		Tp:  1, // 1-广播聊天消息
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	players := WorldMgrObj.GetAllPlayers()

	// 广播给周围的玩家（包括自己）
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

// 同步玩家上线的位置
func (p *Player) SyncSurrounding() {
	// 1 位置msg
	myPosMsg := &pb.BroadCast{
		Pid: p.PID,
		Tp:  2, // 2-标识玩家位置信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 得到当前玩家周边的九宫格信息（包括自己）
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)

	// 2 给周围玩家同步自己的上线位置（包括自己）
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		player := WorldMgrObj.GetPlayerByPid(int32(pid))
		players = append(players, player)
	}
	for _, player := range players {
		player.SendMsg(200, myPosMsg)
	}

	// 3 给自己同步周围玩家的位置信息 --- 202消息
	// 3.1 组建自己周围玩家的位置信息
	surroundingPlayersMsg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			Pid: player.PID,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		surroundingPlayersMsg = append(surroundingPlayersMsg, p)
	}

	// 3.2 组建MsgID:202消息
	surroundingPlayersMsgs := &pb.SyncPlayers{
		Ps: surroundingPlayersMsg,
	}

	// 3.3 发送消息给客户端
	p.SendMsg(202, surroundingPlayersMsgs)
}

// 获取周围的玩家
func (p *Player) GetSurroundingPlayers() (players []*Player) {
	// 得到当前玩家周边的九宫格信息（包括自己）
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)

	// 获取周围玩家
	for _, pid := range pids {
		player := WorldMgrObj.GetPlayerByPid(int32(pid))
		players = append(players, player)
	}
	return
}

// 更新玩家的坐标
func (p *Player) UpdatePosition(x, y, z, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	// 位置移动信息
	msg := &pb.BroadCast{
		Pid: p.PID,
		Tp:  4, // 4-移动之后的位置信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 广播给周围的玩家（包括自己）
	players := p.GetSurroundingPlayers()
	for _, player := range players {
		player.SendMsg(200, msg) // 200表示行为信息
	}
}

// 玩家下线
func (p *Player) Offline() {
	// 得到当前玩家周边的九宫格信息（包括自己）
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)

	// 通知周围玩家自己下线
	msg := &pb.SyncPid{
		Pid: p.PID,
	}

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		player := WorldMgrObj.GetPlayerByPid(int32(pid))
		players = append(players, player)
	}

	for _, player := range players {
		player.SendMsg(201, msg) // 201表示玩家下线
	}
}
