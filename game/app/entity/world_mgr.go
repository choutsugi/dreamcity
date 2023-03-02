package entity

import (
	"dreamcity/shared/code"
	"fmt"
	"github.com/dobyte/due/errors"
)

// WorldMgr 世界管理器
type WorldMgr struct {
	worlds map[string]*World
}

func NewWorldMgr(opts []*WorldOpts) *WorldMgr {
	mgr := &WorldMgr{worlds: make(map[string]*World, len(opts))}
	for i := range opts {
		w := newWorld(opts[i])
		mgr.worlds[w.ID] = w
	}
	return mgr
}

func (mgr *WorldMgr) GetWorld(wid string) (*World, error) {
	scene, ok := mgr.worlds[wid]
	if !ok {
		return nil, errors.NewError(code.NotFoundWorld, fmt.Sprintf("world not exist, wid=%s", wid))
	}
	return scene, nil
}
