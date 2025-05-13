package svc

import (
	"go-microservices/com"
	"go-microservices/service/demoService/internal/config"
)

type ServiceContext struct {
	C      *com.Config
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {

	c.InitLog()
	return &ServiceContext{
		C:      c.InitCom(),
		Config: c,
	}
}
