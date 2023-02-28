package service_test

import (
	"dreamcity/shared/service"
	"github.com/dobyte/due/cluster/node"
	"testing"
)

func TestTokenLogin(t *testing.T) {

	svc := service.NewLogin(&node.Proxy{})
	uid, err := svc.TokenLogin("token", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(uid)
}
