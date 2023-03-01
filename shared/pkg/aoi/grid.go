package aoi

import (
	"fmt"
	"sync"
)

type Grid struct {
	ID   int
	MinX int
	MaxX int
	MinY int
	MaxY int
	pids map[int64]bool
	rw   sync.RWMutex
}

func newGrid(gid, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		ID:   gid,
		MinX: minX,
		MaxX: maxX,
		MinY: minY,
		MaxY: maxY,
		pids: make(map[int64]bool),
	}
}

func (g *Grid) Add(pid int64) {
	g.rw.Lock()
	defer g.rw.Unlock()

	g.pids[pid] = true
}

func (g *Grid) Remove(pid int64) {
	g.rw.Lock()
	defer g.rw.Unlock()

	delete(g.pids, pid)
}

func (g *Grid) GetPids() (pids []int64) {
	g.rw.RLock()
	defer g.rw.RUnlock()

	for pid := range g.pids {
		pids = append(pids, pid)
	}

	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid ID: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, pids:%v",
		g.ID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.pids)
}
