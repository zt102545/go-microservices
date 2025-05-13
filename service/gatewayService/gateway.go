package main

import (
	"flag"
	"fmt"
	_ "github.com/apache/skywalking-go"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"go-microservices/service/gatewayService/internal/config"
	"go-microservices/service/gatewayService/internal/handler"
	"go-microservices/service/gatewayService/internal/svc"
	"os"
	"strings"
)

var configFile = flag.String("f", "./service/gatewayService/etc/gateway.yaml", "the config file")

func main() {

	flag.Parse()
	var c config.Config
	if len(os.Args) > 1 {
		etcd := strings.Split(os.Args[1], ",")
		err := c.NewConfigFromEtcd(etcd, "gatewayConfig")
		if err != nil {
			panic(err)
		}
	} else {
		conf.MustLoad(*configFile, &c)
	}

	s := rest.MustNewServer(c.GetRestConf(), rest.WithCors("*"))
	defer func() {
		s.Stop()
		c.Close()
	}()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(s, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	s.Start()
}
