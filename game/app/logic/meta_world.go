package logic

import (
	"context"
	"dreamcity/game/app/entity"
	"dreamcity/shared/pb/code"
	pb "dreamcity/shared/pb/world"
	"dreamcity/shared/pkg/aoi"
	"dreamcity/shared/pkg/sugar"
	"dreamcity/shared/route"
	"github.com/dobyte/due/cluster"
	"github.com/dobyte/due/cluster/node"
	"github.com/dobyte/due/config"
	"github.com/dobyte/due/log"
	"github.com/dobyte/due/session"
)

type MetaWorld struct {
	proxy     *node.Proxy
	ctx       context.Context
	worldMgr  *entity.WorldMgr
	playerMgr *entity.PlayerMgr
}

func NewMetaWorld(proxy *node.Proxy) *MetaWorld {

	opts := make([]*entity.WorldOpts, 0)
	if err := config.Get("dreamcity.worlds").Scan(&opts); err != nil {
		log.Fatalf("failed to load dreamcity worlds config: %+v\n", err)
	}

	metaWorld := &MetaWorld{
		proxy:     proxy,
		ctx:       context.Background(),
		worldMgr:  entity.NewWorldMgr(opts),
		playerMgr: entity.NewPlayerMgr(),
	}

	return metaWorld
}

func (l *MetaWorld) Init() {
	l.proxy.Events().AddEventHandler(cluster.Disconnect, l.hookDisconnect)
	l.proxy.Router().AddRouteHandler(route.EnterScene, false, l.EnterWorld)
	l.proxy.Router().AddRouteHandler(route.LeaveScene, false, l.LeaveWorld)
}

func (l *MetaWorld) hookDisconnect(event *node.Event) {

	uid := event.UID

	player := l.playerMgr.GetPlayer(uid)
	if player != nil {
		if scene := player.GetWorld(); scene != nil {
			// 获取周围的玩家ID
			targets := scene.AoiMgr.GetPidsByPos(player.PosX, player.PosZ)
			targets = sugar.Delete(targets, uid)
			l.proxy.Multicast(l.ctx, &node.MulticastArgs{
				Kind:    session.User,
				Targets: targets,
				Message: &node.Message{
					Route: route.Broadcast,
					Data: &pb.BroadCast{
						Pid: uid,
						Tp:  pb.BroadCast_PlayerLeave,
					},
				},
			})
			scene.RemPlayer(player)
		}
		l.playerMgr.RemPlayer(player)
	}

	nid := event.Proxy.GetNodeID()
	event.Proxy.UnbindNode(l.ctx, event.UID, nid)
}

func (l *MetaWorld) EnterWorld(ctx *node.Context) {

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
		res.Code = code.Code_Abnormal
		return
	}
	// 检查是否登录
	uid := ctx.Request.UID
	if uid == 0 {
		res.Code = code.Code_NotLogin
		return
	}
	// 玩家ID使用用户ID
	pid := uid
	// 获取场景
	scene, err := l.worldMgr.GetWorld(req.GetSid())
	if err != nil {
		res.Code = code.Code_IllegalParams
		return
	}
	// 检查位置是否正确
	if req.GetPos().GetX() < float32(scene.AoiMgr.MinX) || req.GetPos().GetX() > float32(scene.AoiMgr.MaxX) ||
		req.GetPos().GetZ() < float32(scene.AoiMgr.MinY) || req.GetPos().GetZ() > float32(scene.AoiMgr.MaxY) {
		res.Code = code.Code_IllegalParams
		return
	}
	// 检查玩家是否已在其它场景中
	player := l.playerMgr.GetPlayer(pid)
	if player != nil {
		if scene := player.GetWorld(); scene != nil {
			// 获取周围的玩家ID
			targets := scene.AoiMgr.GetPidsByPos(player.PosX, player.PosZ)
			targets = sugar.Delete(targets, pid)
			l.proxy.Multicast(l.ctx, &node.MulticastArgs{
				Kind:    session.User,
				Targets: targets,
				Message: &node.Message{
					Route: route.Broadcast,
					Data: &pb.BroadCast{
						Pid: pid,
						Tp:  pb.BroadCast_PlayerLeave,
					},
				},
			})
			scene.RemPlayer(player)
		}
	} else {
		player = entity.NewPlayer(pid, req.Pos.X, req.Pos.Y, req.Pos.Z, req.Pos.V)
	}
	// 玩家进入场景
	{
		// 玩家绑定场景
		scene.AddPlayer(player)
		// 玩家管理器添加玩家
		l.playerMgr.AddPlayer(player)
		// 获取周围的玩家ID
		targets := scene.AoiMgr.GetPidsByPos(player.PosX, player.PosZ)
		targets = sugar.Delete(targets, pid)
		// 广播玩家出现
		l.proxy.Multicast(l.ctx, &node.MulticastArgs{
			Kind:    session.User,
			Targets: targets,
			Message: &node.Message{
				Route: route.Broadcast,
				Data: &pb.BroadCast{
					Pid: pid,
					Tp:  pb.BroadCast_PlayerAppear,
					Data: &pb.BroadCast_Player{
						Player: &pb.Player{
							Pid: pid,
							Pos: &pb.Position{
								X: player.PosX,
								Y: player.PosY,
								Z: player.PosZ,
								V: player.PosV,
							},
							Act: &pb.Action{
								Sit:   player.ActSit,
								Jump:  player.ActJump,
								Dance: player.ActDance,
							},
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
				Pid: p.Pid,
				Pos: &pb.Position{
					X: p.PosX,
					Y: p.PosY,
					Z: p.PosZ,
					V: p.PosV,
				},
				Act: &pb.Action{
					Sit:   p.ActSit,
					Jump:  p.ActJump,
					Dance: p.ActDance,
				},
			})
		}
		l.proxy.Push(l.ctx, &node.PushArgs{
			Kind:   session.User,
			Target: pid,
			Message: &node.Message{
				Route: route.SyncArea,
				Data: &pb.SyncArea{
					Ps: surPs,
				},
			},
		})
	}

	ctx.BindNode()

	res.Code = code.Code_Ok
}

func (l *MetaWorld) LeaveWorld(ctx *node.Context) {

	req := &pb.LeaveReq{}
	res := &pb.LeaveRes{}

	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("leave scene response failed, err: %+v\n", err)
		}
	}()

	// 解析参数
	if err := ctx.Request.Parse(req); err != nil {
		log.Errorf("invalid leave_scene message, err: %v", err)
		res.Code = code.Code_Abnormal
		return
	}
	// 检查是否登录
	uid := ctx.Request.UID
	if uid == 0 {
		res.Code = code.Code_NotLogin
		return
	}
	// 玩家ID使用UID
	pid := uid
	// 获取场景
	player := l.playerMgr.GetPlayer(pid)
	if player == nil {
		res.Code = code.Code_NotFound
		return
	}
	scene := player.GetWorld()
	if scene == nil {
		res.Code = code.Code_NotFound
		return
	}
	// 移除玩家
	targets := scene.AoiMgr.GetPidsByPos(player.PosX, player.PosZ)
	targets = sugar.Delete(targets, pid)
	l.proxy.Multicast(l.ctx, &node.MulticastArgs{
		Kind:    session.User,
		Targets: targets,
		Message: &node.Message{
			Route: route.Broadcast,
			Data: &pb.BroadCast{
				Pid: pid,
				Tp:  pb.BroadCast_PlayerLeave,
			},
		},
	})
	scene.RemPlayer(player)

	res.Code = code.Code_Ok
}

func (l *MetaWorld) SyncInfo(ctx *node.Context) {

	req := &pb.SyncInfoReq{}
	res := &pb.SyncInfoRes{}

	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("sync info response failed, err: %+v\n", err)
		}
	}()

	// 解析参数
	if err := ctx.Request.Parse(req); err != nil {
		log.Errorf("invalid sync_info message, err: %v", err)
		res.Code = code.Code_Abnormal
		return
	}
	// 检查是否登录
	uid := ctx.Request.UID
	if uid == 0 {
		res.Code = code.Code_NotLogin
		return
	}
	// 玩家ID使用UID
	pid := uid
	// 获取场景
	player := l.playerMgr.GetPlayer(pid)
	if player == nil {
		res.Code = code.Code_NotFound
		return
	}
	world := player.GetWorld()
	if world == nil {
		res.Code = code.Code_NotFound
		return
	}
	// 检查位置是否正确
	if req.GetPos().GetX() < float32(world.AoiMgr.MinX) || req.GetPos().GetX() > float32(world.AoiMgr.MaxX) ||
		req.GetPos().GetZ() < float32(world.AoiMgr.MinY) || req.GetPos().GetZ() > float32(world.AoiMgr.MaxY) {
		res.Code = code.Code_IllegalParams
		return
	}
	// 獲取移動前後的grid
	oldGid := world.AoiMgr.GetGidByPos(player.PosX, player.PosZ)
	newGid := world.AoiMgr.GetGidByPos(req.GetPos().GetX(), req.GetPos().GetZ())
	// 更新玩家信息
	if req.GetPos() != nil {
		player.PosX = req.GetPos().GetX()
		player.PosY = req.GetPos().GetY()
		player.PosZ = req.GetPos().GetZ()
		player.PosV = req.GetPos().GetV()
	} else {
		player.PosX = 0
		player.PosY = 0
		player.PosZ = 0
		player.PosV = 0
	}
	if req.GetAct() != nil {
		player.ActSit = req.GetAct().GetSit()
		player.ActJump = req.GetAct().GetJump()
		player.ActDance = req.GetAct().GetDance()
	} else {
		player.ActSit = 0
		player.ActJump = 0
		player.ActDance = 0
	}

	if oldGid != newGid {
		world.AoiMgr.RemPidFromGridByGid(pid, oldGid)
		world.AoiMgr.AddPidToGridByGid(pid, newGid)

		oldGrids := world.AoiMgr.GetSurroundGrids(oldGid)
		oldGridsMap := make(map[int]struct{}, len(oldGrids))
		for _, grid := range oldGrids {
			oldGridsMap[grid.ID] = struct{}{}
		}

		newGrids := world.AoiMgr.GetSurroundGrids(newGid)
		newGridsMap := make(map[int]struct{}, len(newGrids))
		for _, grid := range newGrids {
			newGridsMap[grid.ID] = struct{}{}
		}

		// 消失在视野的格子/保持不變的格子
		leavingGrids := make([]*aoi.Grid, 0)
		keepingGrids := make([]*aoi.Grid, 0)
		for _, grid := range oldGrids {
			if _, ok := newGridsMap[grid.ID]; !ok {
				leavingGrids = append(leavingGrids, grid)
			} else {
				keepingGrids = append(keepingGrids, grid)
			}
		}
		// 消失在视野的玩家
		leavingPids := make([]int64, 0)
		for _, grid := range leavingGrids {
			pids := world.AoiMgr.GetPidsByGid(grid.ID)
			leavingPids = append(leavingPids, pids...)
		}
		// 保持不變的玩家
		keepingPids := make([]int64, 0)
		for _, grid := range keepingGrids {
			pids := world.AoiMgr.GetPidsByGid(grid.ID)
			keepingPids = append(keepingPids, pids...)
		}
		// 出現在視野的格子
		enteringGrids := make([]*aoi.Grid, 0)
		for _, grid := range newGrids {
			if _, ok := oldGridsMap[grid.ID]; !ok {
				enteringGrids = append(enteringGrids, grid)
			}
		}
		// 出現在视野的玩家
		enteringPids := make([]int64, 0)
		for _, grid := range enteringGrids {
			pids := world.AoiMgr.GetPidsByGid(grid.ID)
			enteringPids = append(enteringPids, pids...)
		}

	} else {
		l.proxy.Multicast(l.ctx, &node.MulticastArgs{
			Kind:    session.User,
			Targets: world.AoiMgr.GetPidsByPos(player.PosX, player.PosZ),
			Message: &node.Message{
				Route: route.Broadcast,
				Data: &pb.BroadCast{
					Pid: pid,
					Tp:  pb.BroadCast_PlayerInfo,
					Data: &pb.BroadCast_Player{
						Player: &pb.Player{
							Pid: pid,
							Pos: &pb.Position{
								X: player.PosX,
								Y: player.PosY,
								Z: player.PosZ,
								V: player.PosV,
							},
							Act: &pb.Action{
								Sit:   player.ActSit,
								Jump:  player.ActJump,
								Dance: player.ActDance,
							},
						},
					},
				},
			},
		})
	}

	res.Code = code.Code_Ok
}
