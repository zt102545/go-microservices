// Code generated by goctl. DO NOT EDIT.
// Source: demo.proto

package server

import (
	"context"

	"go-microservices/proto/generate/demo"
	"go-microservices/service/demoService/internal/logic"
	"go-microservices/service/demoService/internal/svc"
)

type DemoServer struct {
	svcCtx *svc.ServiceContext
	demo.UnimplementedDemoServer
}

func NewDemoServer(svcCtx *svc.ServiceContext) *DemoServer {
	return &DemoServer{
		svcCtx: svcCtx,
	}
}

// Ping pong
func (s *DemoServer) Ping(ctx context.Context, in *demo.PingReq) (*demo.PingResp, error) {
	l := logic.NewPingLogic(ctx, s.svcCtx)
	return l.Ping(in)
}
