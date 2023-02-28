package main

import (
	"dreamcity/hall/app/route"
	"github.com/dobyte/due"
	"github.com/dobyte/due/cluster/node"
	"github.com/dobyte/due/locate/redis"
	"github.com/dobyte/due/registry/etcd"
	"github.com/dobyte/due/transport/grpc"
)

func main() {
	container := due.NewContainer()
	component := node.NewNode(
		node.WithLocator(redis.NewLocator()),
		node.WithRegistry(etcd.NewRegistry()),
		node.WithTransporter(grpc.NewTransporter()),
	)

	route.Init(component.Proxy())

	container.Add(component)
	container.Serve()
}
