package route

import (
	"dreamcity/hall/app/logic"
	"github.com/dobyte/due/cluster/node"
)

func Init(proxy *node.Proxy) {
	logic.NewLogin(proxy).Init()
}
