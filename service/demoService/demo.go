package main

import (
	"flag"
	"fmt"
	_ "github.com/apache/skywalking-go"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"go-microservices/proto/generate/demo"
	"go-microservices/service/demoService/internal/config"
	"go-microservices/service/demoService/internal/server"
	"go-microservices/service/demoService/internal/svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"os"
	"strings"
)

var configFile = flag.String("f", "./service/demoService/etc/demo.yaml", "the config file")

func main() {

	flag.Parse()
	var c config.Config
	if len(os.Args) > 1 {
		etcd := strings.Split(os.Args[1], ",")
		err := c.NewConfigFromEtcd(etcd, "demoConfig")
		if err != nil {
			panic(err)
		}
	} else {
		conf.MustLoad(*configFile, &c)
	}

	ctx := svc.NewServiceContext(c)
	serverConf := c.GetRpcServerConf()
	s := zrpc.MustNewServer(serverConf, func(grpcServer *grpc.Server) {
		demo.RegisterDemoServer(grpcServer, server.NewDemoServer(ctx))

		if serverConf.Mode == service.DevMode || serverConf.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer func() {
		s.Stop()
		c.Close()
	}()

	fmt.Printf("Starting rpc server at %s:%d...\n", c.Host, c.Port)
	s.Start()
}
