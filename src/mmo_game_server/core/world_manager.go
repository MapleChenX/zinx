package core

import "sync"

// 游戏世界管理模块
type WorldManager struct {
	// 当前世界地图
	AoiMgr *AOIManager
	// 玩家集合
	Players map[int32]*Player
	pLock   sync.RWMutex
}

var WorldMgrObj *WorldManager

// 初始化世界管理模块
func init() {
	WorldMgrObj = &WorldManager{
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_X, AOI_CNTS_Y),
		Players: make(map[int32]*Player),
	}
}

// 添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.PID] = player
	wm.pLock.Unlock()

	// 将player添加到世界地图中
	wm.AoiMgr.AddToGridByPos(int(player.PID), player.X, player.Z)
}

// 删除一个玩家
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	wm.pLock.Lock()
	defer wm.pLock.Unlock()

	player := wm.Players[pid]

	// 将player从世界地图中删除
	wm.AoiMgr.RemoveFromGridByPos(int(player.PID), player.X, player.Z)

	delete(wm.Players, pid)
}

// 通过玩家ID获取玩家对象
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pid]
}

// 获取全部玩家
func (wm *WorldManager) GetAllPlayers() (players []*Player) {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	for _, player := range wm.Players {
		players = append(players, player)
	}

	return
}
