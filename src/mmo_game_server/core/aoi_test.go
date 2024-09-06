package core

import "testing"

func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	t.Logf("aoiMgr: %v\n", aoiMgr)
}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	grids := aoiMgr.GetSurroundGridsByGid(24)
	// 循环打印
	for _, grid := range grids {
		t.Logf(grid.String())
	}
}
