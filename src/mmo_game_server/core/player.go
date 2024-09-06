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

var PidGen int32 = 114514
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

// 发送消息给客户端肉
func (p *Player) SendMsg(msgId uint32, msg proto.Message) {
	if p.Conn == nil {
		fmt.Println("Connection in player is nil")
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
