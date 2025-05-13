package svc

import (
	"context"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"go-microservices/proto/client"
	"go-microservices/service/gatewayService/internal/config"
	"go-microservices/service/gatewayService/internal/middleware"
	uconfig "go-microservices/utils/config"
	"go-microservices/utils/logs"
)

type ServiceContext struct {
	Config          config.Config
	AuthInterceptor rest.Middleware

	RpcDemo client.Demo
}

func NewServiceContext(c config.Config) *ServiceContext {

	c.InitLog()

	demoClient, err := zrpc.NewClient(c.GetRpcClientConf(uconfig.DemoService))
	if err != nil {
		logs.Err(context.Background(), "RpcDemo client init error: %v", err, logs.Flag("Init"))
	}

	return &ServiceContext{
		Config:          c,
		AuthInterceptor: middleware.NewAuthInterceptorMiddleware().Handle,

		RpcDemo: client.NewDemo(demoClient),
	}
}
