package logic

import (
	"context"
	"dreamcity/shared/pb/code"
	pb "dreamcity/shared/pb/login"
	"dreamcity/shared/route"
	"dreamcity/shared/service"
	"github.com/dobyte/due/cluster/node"
	"github.com/dobyte/due/log"
)

type Login struct {
	proxy    *node.Proxy
	ctx      context.Context
	loginSvc *service.Login
}

func NewLogin(proxy *node.Proxy) *Login {
	return &Login{
		proxy:    proxy,
		ctx:      context.Background(),
		loginSvc: service.NewLogin(proxy),
	}
}

// 登录请求
func (l *Login) login(ctx *node.Context) {

	req := &pb.LoginReq{}
	res := &pb.LoginRes{}

	// 响应
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("login response failed, err: %+v\n", err)
		}
	}()

	// 解析请求
	if err := ctx.Request.Parse(req); err != nil {
		log.Errorf("invalid login message, err: %v", err)
		res.Code = code.Code_Abnormal
		return
	}

	/*
		// 获取IP
		ip, err := ctx.GetIP()
		if err != nil {
			log.Errorf("get client ip failed, err: %v", err)
			res.State = false
			return
		}

		// 登录逻辑
		uid, err := l.loginSvc.TokenLogin(req.GetToken(), ip)
		if err != nil {
			log.Errorf("login failed, err: %v", err)
			res.State = false
			return
		}
	*/

	// 绑定网关
	if err := ctx.BindGate(req.Uid); err != nil {
		log.Errorf("bind gate failed, err: %v", err)
		res.Code = code.Code_Failed
		return
	}

	// 响应
	res.Code = code.Code_Ok
	log.Infof("登录成功, uid = %d\n", req.Uid)
}

func (l *Login) Init() {
	// 注册路由
	l.proxy.Router().AddRouteHandler(route.Login, false, l.login)
}
