package entity

import (
	"dreamcity/shared/pkg/aoi"
	"sync"
)

type Scene struct {
	ID        string
	Name      string
	GridMgr   *aoi.GridMgr
	PlayerNum int
	players   map[int64]*Player
	rw        sync.RWMutex
}

func newScene(opt *SceneOpts) *Scene {
	gridMgr := aoi.NewGridMgr(opt.MinX, opt.MaxX, opt.ContsX, opt.MinY, opt.MaxY, opt.ContsY)
	return &Scene{
		ID:      opt.ID,
		Name:    opt.Name,
		GridMgr: gridMgr,
		players: make(map[int64]*Player),
	}
}

// AddPlayer 添加玩家
func (s *Scene) AddPlayer(player *Player) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	player.scene = s
	s.players[player.UID()] = player
	s.GridMgr.AddPidToGridByPos(player.UID(), player.PosX, player.PosZ)
	s.PlayerNum++

	return nil
}

// RemPlayer 移除玩家
func (s *Scene) RemPlayer(player *Player) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.GridMgr.RemPidFromGridByPos(player.UID(), player.PosX, player.PosZ)
	delete(s.players, player.UID())
	player.scene = nil
	s.PlayerNum--

	return nil
}
