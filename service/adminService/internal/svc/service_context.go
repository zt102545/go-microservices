package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"go-microservices/com"
	"go-microservices/service/adminService/internal/config"
	"go-microservices/service/adminService/internal/middleware"
)

type ServiceContext struct {
	C               *com.Config
	Config          config.Config
	AuthInterceptor rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {

	c.InitLog()
	return &ServiceContext{
		C:               c.InitCom(),
		Config:          c,
		AuthInterceptor: middleware.NewAuthInterceptorMiddleware().Handle,
	}
}
