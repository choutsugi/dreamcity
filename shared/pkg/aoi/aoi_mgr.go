package aoi

import "fmt"

type Mgr struct {
	MinX  int
	MaxX  int
	CntsX int
	MinY  int
	MaxY  int
	CntsY int
	grids map[int]*Grid
}

func NewAoiMgr(minX, maxX, cntsX, minY, maxY, cntsY int) *Mgr {
	mgr := &Mgr{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			gid := y*cntsX + x
			mgr.grids[gid] = newGrid(gid,
				mgr.MinX+x*mgr.gridWidth(),
				mgr.MinX+(x+1)*mgr.gridWidth(),
				mgr.MinY+y*mgr.gridLength(),
				mgr.MinY+(y+1)*mgr.gridLength())
		}
	}

	return mgr
}

func (m *Mgr) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

func (m *Mgr) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

func (m *Mgr) String() string {
	str := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n GrIDs in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		str += fmt.Sprintln(grid)
	}
	return str
}

func (m *Mgr) GetSurroundGrids(gid int) []*Grid {

	if _, ok := m.grids[gid]; !ok {
		return nil
	}

	grids := make([]*Grid, 0)
	grids = append(grids, m.grids[gid])

	x, y := gid%m.CntsX, gid/m.CntsX
	gids := make([]int, 0)
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	for i := 0; i < 8; i++ {
		newX := x + dx[i]
		newY := y + dy[i]
		if newX >= 0 && newX < m.CntsX && newY >= 0 && newY < m.CntsY {
			gids = append(gids, newY*m.CntsX+newX)
		}
	}

	for _, id := range gids {
		grids = append(grids, m.grids[id])
	}

	return grids
}

func (m *Mgr) GetGidByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(y) - m.MinY) / m.gridLength()

	return gy*m.CntsX + gx
}

func (m *Mgr) GetPidsByPos(x, y float32) []int64 {
	pids := make([]int64, 0)
	gid := m.GetGidByPos(x, y)
	grids := m.GetSurroundGrids(gid)
	for _, v := range grids {
		pids = append(pids, v.GetPids()...)
	}

	return pids
}

func (m *Mgr) GetPidsByGid(gid int) []int64 {
	pids := make([]int64, 0)
	pids = m.grids[gid].GetPids()
	return pids
}

func (m *Mgr) RemPidFromGridByGid(pid int64, gid int) {
	m.grids[gid].Remove(pid)
}

func (m *Mgr) AddPidToGridByGid(pid int64, gid int) {
	m.grids[gid].Add(pid)
}

func (m *Mgr) AddPidToGridByPos(pid int64, x, y float32) {
	gid := m.GetGidByPos(x, y)
	grid := m.grids[gid]
	grid.Add(pid)
}

func (m *Mgr) RemPidFromGridByPos(pid int64, x, y float32) {
	gid := m.GetGidByPos(x, y)
	grid := m.grids[gid]
	grid.Remove(pid)
}
