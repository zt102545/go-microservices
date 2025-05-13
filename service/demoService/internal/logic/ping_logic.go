package logic

import (
	"context"
	"fmt"
	"go-microservices/utils/logs"

	"go-microservices/proto/generate/demo"
	"go-microservices/service/demoService/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Ping pong
func (l *PingLogic) Ping(in *demo.PingReq) (*demo.PingResp, error) {

	//redis操作示例
	err := l.svcCtx.C.LockRedis.Lock(l.ctx, "ping", 60)
	if err != nil {
		logs.Err(l.ctx, "lock is err :%s", err, logs.Flag("lock"))
		return &demo.PingResp{}, err
	}

	//mysql操作示例
	user, err := l.svcCtx.C.UserModel.FindOne(l.ctx, 1)
	if err != nil {
		logs.Err(l.ctx, "UserModel is err:%v", err, logs.Flag("Ping"))
		return &demo.PingResp{Pong: err.Error()}, nil
	}

	//公共接口操作示例
	globalVariables, err := l.svcCtx.C.GetGlobalVariable(l.ctx, 1)
	if err != nil {
		logs.Err(l.ctx, "GlobalVariablesModel is err:%v", err, logs.Flag("Ping"))
		return &demo.PingResp{Pong: err.Error()}, nil
	}
	//测试日志
	logs.Info(l.ctx, fmt.Sprintf("PingLogic PingReq: %v", in))
	return &demo.PingResp{Pong: user.Username + globalVariables.Name}, nil
}
