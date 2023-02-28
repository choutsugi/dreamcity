package main

import (
	"github.com/dobyte/due"
	"github.com/dobyte/due/cluster/gate"
	"github.com/dobyte/due/locate/redis"
	"github.com/dobyte/due/network/tcp"
	"github.com/dobyte/due/registry/etcd"
	"github.com/dobyte/due/transport/grpc"
)

func main() {
	container := due.NewContainer() // 创建容器
	component := gate.NewGate(
		gate.WithServer(tcp.NewServer()),            // tcp
		gate.WithLocator(redis.NewLocator()),        // redis
		gate.WithRegistry(etcd.NewRegistry()),       // etcd
		gate.WithTransporter(grpc.NewTransporter()), // grpc
	) // 创建组件

	container.Add(component) // 添加组件
	container.Serve()        // 运行容器
}
