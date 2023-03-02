package entity

import "sync"

var (
	playerMgrOnce     sync.Once
	playerMgrInstance *PlayerMgr
)

type PlayerMgr struct {
	// TODO：UserService
	rw      sync.RWMutex
	players map[int64]*Player
}

func NewPlayerMgr() *PlayerMgr {
	playerMgrOnce.Do(func() {
		playerMgrInstance = &PlayerMgr{
			players: make(map[int64]*Player),
		}
	})
	return playerMgrInstance
}

func (mgr *PlayerMgr) AddPlayer(player *Player) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	mgr.players[player.Pid] = player
}

func (mgr *PlayerMgr) RemPlayer(player *Player) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	delete(mgr.players, player.Pid)
}

func (mgr *PlayerMgr) GetPlayer(pid int64) *Player {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	return mgr.players[pid]
}
