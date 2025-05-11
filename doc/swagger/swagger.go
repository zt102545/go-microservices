package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"go-microservices/doc/swagger/internal/config"
	"go-microservices/doc/swagger/internal/handler"
)

var configFile = flag.String("f", "./doc/swagger/etc/swagger-api.yaml", "the config file")

// 生成的swagger有个坑点，content type没法用form-data请求，故暂停使用，生成的json文档可以导入到其他请求工具中使用，如ApiFox
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
