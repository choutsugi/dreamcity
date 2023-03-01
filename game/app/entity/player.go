package entity

import (
	"dreamcity/shared/model/user"
	"sync"
)

type Player struct {
	user  *user.User
	rw    sync.RWMutex
	scene *Scene
	PosX  float32
	PosY  float32
	PosZ  float32
	PosV  float32
}

func NewPlayer(user *user.User, posX, posY, posZ, posV float32) *Player {
	return &Player{
		user: user,
		PosX: posX,
		PosY: posY,
		PosZ: posZ,
		PosV: posV,
	}
}

func (p *Player) UID() int64 {
	return p.user.UID
}

func (p *Player) User() *user.User {
	u := p.user
	return u
}

func (p *Player) GetScene() *Scene {
	scene := p.scene
	return scene
}
