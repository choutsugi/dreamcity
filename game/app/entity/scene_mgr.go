package entity

import "github.com/dobyte/due/errors"

// SceneMgr 场景管理器
type SceneMgr struct {
	scenes map[string]*Scene
}

func NewSceneMgr(opts []*SceneOpts) *SceneMgr {
	mgr := &SceneMgr{scenes: make(map[string]*Scene, len(opts))}
	for i := range opts {
		scene := newScene(opts[i])
		mgr.scenes[scene.ID] = scene
	}
	return mgr
}

// GetScene 获取场景
func (mgr *SceneMgr) GetScene(sceneID string) (*Scene, error) {
	scene, ok := mgr.scenes[sceneID]
	if !ok {
		return nil, errors.New("场景不存在")
	}
	return scene, nil
}
