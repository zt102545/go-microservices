package config

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"go-microservices/com"
	"go-microservices/utils/kafka"
	"go-microservices/utils/logs"
	"go-microservices/utils/mysql"
	"go-microservices/utils/postgres"
	"go-microservices/utils/redis"
	"go.etcd.io/etcd/client/v3"
	"time"
)

type Config struct {
	Name string `json:"name"`
	Env  string `json:"env,default=pro"`
	Host string `json:"host,default=0.0.0.0"`
	Port int    `json:"port,default=0"`
	Etcd struct {
		Hosts []string `json:"hosts,optional"`
		Key   string   `json:"key,optional"`
	} `json:"etcd,optional"`
	Db       mysql.DbConfig          `json:"db,optional"`
	Postgres postgres.PostgresConfig `json:"postgres,optional"`
	Redis    redis.RedisConfig       `json:"redis,optional"`
	Kafka    kafka.KafkaConfig       `json:"elastic,optional"`
	Log      logs.LoggerConfig       `json:"resources,optional"`
}

const (
	DemoService = "DemoRpc"

	DemoConfig = "DemoConfig"
)

var ServiceConfig *Config

// GetRestConf 获取api配置
func (c *Config) GetRestConf() rest.RestConf {
	restConf := rest.RestConf{
		ServiceConf: service.ServiceConf{
			Name: c.Name,
		},
		Host: c.Host,
		Port: c.Port,
	}
	return restConf
}

// GetRpcServerConf 获取PRC服务端配置
func (c *Config) GetRpcServerConf() zrpc.RpcServerConf {
	listenOn := fmt.Sprintf("%s:%d", c.Host, c.Port)
	rpcServerConf := zrpc.RpcServerConf{
		ServiceConf: service.ServiceConf{
			Name: c.Name,
			Mode: c.Env,
		},
		ListenOn: listenOn,
		Etcd: discov.EtcdConf{
			Hosts: c.Etcd.Hosts,
			Key:   c.Etcd.Key,
		},
	}
	return rpcServerConf
}

// GetRpcClientConf 获取RPC客户端配置
func (c *Config) GetRpcClientConf(serviceName string) zrpc.RpcClientConf {
	rpcClientConf := zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: c.Etcd.Hosts,
			Key:   serviceName,
		},
	}
	return rpcClientConf
}

func (c *Config) InitLog() *logs.Logger {
	//修改go-zero框架自带日志的输出模式
	//logx.SetWriter(logx.NewWriter(kafka.NewProducer(c.Log.KafkaInfo.Address, c.Log.KafkaInfo.Topic)))
	c.Log.Env = c.Env
	return logs.NewLogger(c.Log)
}

// InitCom 初始化com配置
func (c *Config) InitCom() *com.Config {

	// 全局配置
	ServiceConfig = c
	// 初始化日志
	c.InitLog()
	// 初始化Kafka生产者
	c.Kafka.InitProducer()
	// 初始化Kafka消费组
	c.Kafka.InitConsumer()
	// 初始化Com
	return com.InitCom(c.Db.Init(), c.Postgres.Init(), c.Redis.Init())
}

// EtcdChangeWatch 动态监听etcd变化，实时更新数据
func (c *Config) EtcdChangeWatch(addr []string, key string, f func(value []byte)) {

	ctx := context.Background()
	defer func() {
		if err := recover(); err != nil {
			logs.Err(ctx, "etcd failed, err:%v ", err, logs.Flag("etcd"))
		}
	}()

	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 30 * time.Second,
	})
	if err != nil {
		logs.Err(ctx, fmt.Sprintf("connect to etcd failed, err:%v ", err), logs.Flag("etcd"))
		return
	}
	defer etcdCli.Close()

	// watch key change
	rch := etcdCli.Watch(context.Background(), key) // <-chan WatchResponse
	for resp := range rch {
		if len(resp.Events) == 0 {
			continue
		}
		logs.Info(ctx, fmt.Sprintf("etcd change %s", resp.Events[0].Kv.Value), logs.Flag("etcd"))
		// 更新配置
		f(resp.Events[0].Kv.Value)
	}
}

func (c *Config) Close() {
	c.Db.Close()
	c.Redis.Close()
}
