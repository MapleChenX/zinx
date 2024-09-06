package core

import "fmt"

// 定义边界值
const (
	AOI_MIN_X  = 85
	AOI_MAX_X  = 410
	AOI_CNTS_X = 10
	AOI_MIN_Y  = 75
	AOI_MAX_Y  = 400
	AOI_CNTS_Y = 20
)

/**
 * AOI区域管理模块 --- 地图中的一个区域
 */

type AOIManager struct {
	// 区域的左边界坐标
	MinX int
	// 区域的右边界坐标
	MaxX int
	// X方向格子的数量
	CntsX int
	// 区域的上边界坐标
	MinY int
	// 区域的下边界坐标
	MaxY int
	// Y方向格子的数量
	CntsY int
	// 当前区域中有哪些格子，key=格子ID，value=格子对象
	grids map[int]*Grid
}

// 初始化一个AOI区域管理模块
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	// 给AOI初始化区域的格子所有的格子进行编号和初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算格子ID
			gid := y*cntsX + x
			// 初始化一个格子对象
			aoiMgr.grids[gid] = NewGrid(gid,
				minX+x*aoiMgr.gridWidth(),
				minX+(x+1)*aoiMgr.gridWidth(),
				minY+y*aoiMgr.gridLength(),
				minY+(y+1)*aoiMgr.gridLength())
		}
	}

	return aoiMgr
}

// 得到每个格子在X轴方向的宽度
func (a *AOIManager) gridWidth() int {
	return (a.MaxX - a.MinX) / a.CntsX
}

// 得到每个格子在Y轴方向的高度
func (a *AOIManager) gridLength() int {
	return (a.MaxY - a.MinY) / a.CntsY
}

// 打印出AOIManager信息
func (a *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX: %d, MaxX: %d, CntsX: %d, MinY: %d, MaxY: %d, CntsY: %d\n Grids in AOIManager:\n",
		a.MinX, a.MaxX, a.CntsX, a.MinY, a.MaxY, a.CntsY)
	for _, grid := range a.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据格子的GID得到周边九宫格格子的ID集合
func (a *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	// 判断GID是否在AOIManager中
	if _, ok := a.grids[gID]; !ok {
		return
	}

	// 将当前gid添加到九宫格中
	grids = append(grids, a.grids[gID])

	// 根据GID得到当前格子的X轴编号 idx = id % CntsX
	idx := gID % a.CntsX
	// 判断idx左边是否有格子，如果有，加入到九宫格中
	if idx > 0 {
		grids = append(grids, a.grids[gID-1])
	}
	// 判断idx右边是否有格子，如果有，加入到九宫格中
	if idx < a.CntsX-1 {
		grids = append(grids, a.grids[gID+1])
	}

	// 将X轴当前的格子都取出，进行遍历，再分别得到Y轴上下是否有格子
	for _, grid := range grids {
		// 得到当前格子的ID
		gID := grid.GID
		// 得到当前格子属于哪一行
		idy := gID / a.CntsX
		// 判断当前行的上一行是否有格子，如果有，加入到九宫格中
		if idy > 0 {
			grids = append(grids, a.grids[gID-a.CntsX])
		}
		// 判断当前行的下一行是否有格子，如果有，加入到九宫格中
		if idy < a.CntsY-1 {
			grids = append(grids, a.grids[gID+a.CntsX])
		}
	}

	// 到这里，grids保存了九宫格的所有格子
	return
}

// 通过横纵坐标得到当前坐标属于哪个格子ID
func (a *AOIManager) GetGIDByPos(x, y float32) int {
	gID := (int(y)-a.MinY)/a.gridLength()*a.CntsX + (int(x)-a.MinX)/a.gridWidth()
	return gID
}

// 根据格子的坐标得到周边九宫格格子的ID集合
func (a *AOIManager) GetSurroundGridsByPos(x, y float32) (grids []*Grid) {
	// 得到当前坐标属于哪个格子ID
	gID := a.GetGIDByPos(x, y)
	// 根据格子的GID得到周边九宫格格子的ID集合
	grids = a.GetSurroundGridsByGid(gID)
	return
}

// 根据格子坐标获取周围九宫格内全部的PlayerIDs
func (a *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	// 根据格子ID得到周边九宫格的格子集合
	grids := a.GetSurroundGridsByPos(x, y)
	// 将九宫格内的所有Player的ID累加到playerIDs
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
	}
	return
}

// 添加一个PlayerID到一个格子中
func (a *AOIManager) AddToGridByPos(playerID int, x, y float32) {
	// 根据坐标得到当前属于哪个格子
	gID := a.GetGIDByPos(x, y)
	// 根据格子ID得到格子对象
	grid := a.grids[gID]
	// 将PlayerID添加到格子中
	grid.Add(playerID)
}

// 移除一个PlayerID从一个格子中
func (a *AOIManager) RemoveFromGridByPos(playerID int, x, y float32) {
	// 根据坐标得到当前属于哪个格子
	gID := a.GetGIDByPos(x, y)
	// 根据格子ID得到格子对象
	grid := a.grids[gID]
	// 将PlayerID从格子中删除
	grid.Remove(playerID)
}

// 通过GID获取该格子下全部的PlayerIDs
func (a *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	// 判断GID是否在AOIManager中
	if _, ok := a.grids[gID]; !ok {
		return
	}
	// 得到当前格子的所有PlayerID
	playerIDs = a.grids[gID].GetPlayerIDs()
	return
}

// 添加一个PlayerID到一个格子中
func (a *AOIManager) AddToGridByGid(playerID, gID int) {
	// 判断GID是否在AOIManager中
	if _, ok := a.grids[gID]; !ok {
		return
	}
	// 根据格子ID得到格子对象
	grid := a.grids[gID]
	// 将PlayerID添加到格子中
	grid.Add(playerID)
}

// 移除一个PlayerID从一个格子中
func (a *AOIManager) RemoveFromGridByGid(playerID, gID int) {
	// 判断GID是否在AOIManager中
	if _, ok := a.grids[gID]; !ok {
		return
	}
	// 根据格子ID得到格子对象
	grid := a.grids[gID]
	// 将PlayerID从格子中删除
	grid.Remove(playerID)
}
