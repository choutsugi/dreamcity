package route

import (
	"dreamcity/game/app/logic"
	"github.com/dobyte/due/cluster/node"
)

func Init(proxy *node.Proxy) {
	logic.NewMetaWorld(proxy).Init()
}
