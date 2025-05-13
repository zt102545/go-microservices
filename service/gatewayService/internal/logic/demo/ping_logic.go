package demo

import (
	"context"
	"go-microservices/proto/generate/demo"

	"go-microservices/service/gatewayService/internal/svc"
	"go-microservices/service/gatewayService/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// ping
func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PingLogic) Ping(req *types.PingReq) (resp *types.PingResp, err error) {

	res, err := l.svcCtx.RpcDemo.Ping(l.ctx, &demo.PingReq{
		Ping: req.Ping,
	})
	if err != nil {
		return &types.PingResp{}, err
	}

	resp = &types.PingResp{
		Pong: res.Pong,
	}
	return
}
