package entity

type Player struct {
	Pid      int64
	world    *World
	PosX     float32
	PosY     float32
	PosZ     float32
	PosV     float32
	ActSit   int32
	ActJump  int32
	ActDance int32
}

func NewPlayer(pid int64, posX, posY, posZ, posV float32) *Player {
	return &Player{
		Pid:  pid,
		PosX: posX,
		PosY: posY,
		PosZ: posZ,
		PosV: posV,
	}
}

func (p *Player) GetWorld() *World {
	return p.world
}
