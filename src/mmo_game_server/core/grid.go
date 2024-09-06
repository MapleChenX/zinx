package core

import (
	"fmt"
	"sync"
)

/**
 * 区域中的一个方格
 */

type Grid struct {
	// 网格ID
	GID int
	// 网格左边界
	MinX int
	// 网格右边界
	MaxX int
	// 网格上边界
	MinY int
	// 网格下边界
	MaxY int

	// 网格玩家集合
	PlayerIds  map[int]bool
	PlayerLock sync.RWMutex
}

// 初始化网格
func NewGrid(gid, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gid,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		PlayerIds: make(map[int]bool),
	}
}

// 添加一个玩家
func (g *Grid) Add(playerID int) {
	g.PlayerLock.Lock()
	defer g.PlayerLock.Unlock()
	g.PlayerIds[playerID] = true
}

// 删除一个玩家
func (g *Grid) Remove(playerID int) {
	g.PlayerLock.Lock()
	defer g.PlayerLock.Unlock()
	delete(g.PlayerIds, playerID)

}

// 获取当前网格所有玩家id
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.PlayerLock.RLock()
	defer g.PlayerLock.RUnlock()

	for playerID := range g.PlayerIds {
		playerIDs = append(playerIDs, playerID)
	}
	return
}

// 获取当前网格所有玩家
func (g *Grid) GetPlayers() (players []*Player) {
	g.PlayerLock.RLock()
	defer g.PlayerLock.RUnlock()
	for playerID := range g.PlayerIds {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(playerID)))
	}
	return
}

// 获取当前网格所有玩家数量
func (g *Grid) GetPlayerCount() int {
	g.PlayerLock.RLock()
	defer g.PlayerLock.RUnlock()
	return len(g.PlayerIds)
}

// 打印格子信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.PlayerIds)
}
