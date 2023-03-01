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

// TODO：加载玩家、获取玩家、卸载玩家 => 添加到玩家管理器 => 添加到场景管理器

func (mgr *PlayerMgr) AddPlayer(player *Player) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	mgr.players[player.UID()] = player
}

func (mgr *PlayerMgr) RemPlayer(player *Player) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	delete(mgr.players, player.UID())
}

func (mgr *PlayerMgr) GetPlayer(uid int64) *Player {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	return mgr.players[uid]
}
