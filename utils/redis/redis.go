package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisConfig struct {
	Addr     []string `json:"addr"`
	Password string   `json:"password,optional"`
}

var rdb *redis.Client

func GetClient() *redis.Client {
	return rdb
}

func (r *RedisConfig) Init() *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:         r.Addr[0],        // 地址。
		Password:     r.Password,       // 密码
		DialTimeout:  10 * time.Second, // 设置连接超时
		ReadTimeout:  30 * time.Second, // 设置读取超时
		WriteTimeout: 30 * time.Second, // 设置写入超时
		PoolSize:     500,              // 连接池最大socket连接数，默认为5倍CPU数， 5 * runtime.NumCPU
		MinIdleConns: 10,               // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		PoolTimeout:  30 * time.Second, // 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。
		MaxRetries:   2,                // 命令执行失败时，最多重试多少次，默认为0即不重试
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong, err)
	return rdb
}

func (r *RedisConfig) Close() {
	if rdb != nil {
		err := rdb.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
