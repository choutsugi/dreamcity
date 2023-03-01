package service

import (
	"context"
	"github.com/dobyte/due/cluster/node"
	"github.com/dobyte/due/log"
	"sync/atomic"
)

var (
	UID = int64(100000)
)

type Login struct {
	ctx   context.Context
	proxy *node.Proxy
}

func NewLogin(proxy *node.Proxy) *Login {
	return &Login{
		ctx:   context.Background(),
		proxy: proxy,
	}
}

func (svc *Login) TokenLogin(token string, clientIP string) (int64, error) {
	// TODO：解析Token

	// TODO：持久化登录记录
	log.Infof("user login by token, token=%s, clientIP=%s\n", token, clientIP)

	return atomic.AddInt64(&UID, 1), nil
}
