package entity

import (
	"dreamcity/shared/pkg/aoi"
	"sync"
)

type World struct {
	ID        string
	Name      string
	AoiMgr    *aoi.Mgr
	PlayerNum int
	players   map[int64]*Player
	pLock     sync.RWMutex
}

func newWorld(opt *WorldOpts) *World {
	gridMgr := aoi.NewAoiMgr(opt.MinX, opt.MaxX, opt.ContsX, opt.MinY, opt.MaxY, opt.ContsY)
	return &World{
		ID:      opt.ID,
		Name:    opt.Name,
		AoiMgr:  gridMgr,
		players: make(map[int64]*Player),
	}
}

func (s *World) AddPlayer(player *Player) {
	s.pLock.Lock()
	defer s.pLock.Unlock()

	player.world = s
	s.players[player.Pid] = player
	s.AoiMgr.AddPidToGridByPos(player.Pid, player.PosX, player.PosZ)
	s.PlayerNum++
}

func (s *World) RemPlayer(player *Player) {
	s.pLock.Lock()
	defer s.pLock.Unlock()

	s.AoiMgr.RemPidFromGridByPos(player.Pid, player.PosX, player.PosZ)
	delete(s.players, player.Pid)
	player.world = nil
	s.PlayerNum--
}
