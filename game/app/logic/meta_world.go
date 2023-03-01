package logic

import (
	"context"
	"dreamcity/game/app/entity"
	"dreamcity/shared/model/user"
	pb "dreamcity/shared/pb/scene"
	"github.com/dobyte/due/cluster"
	"github.com/dobyte/due/cluster/node"
	"github.com/dobyte/due/config"
	"github.com/dobyte/due/log"
	"github.com/dobyte/due/session"
)

type MetaWorld struct {
	proxy     *node.Proxy
	ctx       context.Context
	sceneMgr  *entity.SceneMgr
	playerMgr *entity.PlayerMgr
}

func NewMetaWorld(proxy *node.Proxy) *MetaWorld {

	opts := make([]*entity.SceneOptions, 0)
	if err := config.Get("dreamcity.scenes").Scan(&opts); err != nil {
		log.Fatalf("failed to load dreamcity scenes config: %+v\n", err)
	}

	return &MetaWorld{
		proxy:     proxy,
		ctx:       context.Background(),
		sceneMgr:  entity.NewSceneMgr(opts), // 场景管理器
		playerMgr: entity.NewPlayerMgr(),    // 玩家管理器
	}
}

func (l *MetaWorld) Init() {
	l.proxy.Events().AddEventHandler(cluster.Disconnect, l.disconnect)
	l.proxy.Router().AddRouteHandler(2, false, l.enterScene)
}

func (l *MetaWorld) disconnect(event *node.Event) {
	player := l.playerMgr.GetPlayer(event.UID)
	if player != nil {
		if scene := player.GetScene(); scene != nil {
			scene.RemPlayer(player)
		}
		l.playerMgr.RemPlayer(player)
	}

	nid := event.Proxy.GetNodeID()
	event.Proxy.UnbindNode(l.ctx, event.UID, nid)
}

func (l *MetaWorld) enterScene(ctx *node.Context) {

	req := &pb.EnterReq{}
	res := &pb.EnterRes{}

	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("enter scene response failed, err: %+v\n", err)
		}
	}()

	// 解析参数
	if err := ctx.Request.Parse(req); err != nil {
		log.Errorf("invalid enter_scene message, err: %v", err)
		res.Msg = "消息未注册"
		return
	}
	// 检查是否登录
	uid := ctx.Request.UID
	if uid == 0 {
		res.Msg = "未登录"
		return
	}
	// 获取场景
	scene, err := l.sceneMgr.GetScene(req.GetSid())
	if err != nil {
		res.Msg = "场景不存在"
		return
	}
	// 检查位置是否正确
	if req.GetPos().GetX() < float32(scene.GridMgr.MinX) || req.GetPos().GetX() > float32(scene.GridMgr.MaxX) ||
		req.GetPos().GetZ() < float32(scene.GridMgr.MinY) || req.GetPos().GetZ() > float32(scene.GridMgr.MaxY) {
		res.Msg = "请求参数错误"
		return
	}
	// 检查玩家是否已在其它场景中
	player := l.playerMgr.GetPlayer(uid)
	if player != nil {
		if scene := player.GetScene(); scene != nil {
			// todo：广播玩家离开
			scene.RemPlayer(player)
			player.SetScene(nil)
		}
	} else {
		player = entity.NewPlayer(&user.User{UID: uid}, req.Pos.X, req.Pos.Y, req.Pos.Z, req.Pos.V)
	}
	// 玩家进入场景
	{
		// 玩家绑定场景
		player.SetScene(scene)
		scene.AddPlayer(player)
		// 玩家ID添加到Grid
		scene.GridMgr.AddPidToGridByPos(player.UID(), player.PosX, player.PosZ)
		// 玩家管理器添加玩家
		l.playerMgr.AddPlayer(player)
		// 获取周围的玩家ID
		targets := scene.GridMgr.GetPidsByPos(player.PosX, player.PosZ)
		// 广播玩家出现
		l.proxy.Multicast(l.ctx, &node.MulticastArgs{
			Kind:    session.User,
			Targets: targets,
			Message: &node.Message{
				Route: 3,
				Data: &pb.BroadCast{
					Pid: uid,
					Tp:  pb.BroadCast_PlayerAppear,
					Data: &pb.BroadCast_Pos{
						Pos: &pb.Position{
							X: player.PosX,
							Y: player.PosY,
							Z: player.PosZ,
							V: player.PosV,
						},
					},
				},
			},
		})
		// 推送周围玩家信息
		surPs := make([]*pb.Player, 0, len(targets))
		for _, target := range targets {
			p := l.playerMgr.GetPlayer(target)
			surPs = append(surPs, &pb.Player{
				Pid: p.UID(),
				Pos: &pb.Position{
					X: p.PosX,
					Y: p.PosY,
					Z: p.PosZ,
					V: p.PosV,
				},
				Act: nil, // TODO
			})
		}
		l.proxy.Push(l.ctx, &node.PushArgs{
			Kind:   session.User,
			Target: uid,
			Message: &node.Message{
				Route: 4,
				Data: &pb.SyncArea{
					Ps: surPs,
				},
			},
		})
	}

	ctx.BindNode()

	res.State = true
}
