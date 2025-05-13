package config

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"go-microservices/utils/config"
	"go-microservices/utils/gos"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type Config struct {
	config.Config
}

func (c *Config) NewConfigFromEtcd(addr []string, key string) error {

	// 创建etcd客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,             // etcd服务地址
		DialTimeout: 10 * time.Second, // 连接超时时间
	})
	if err != nil {
		log.Fatal("Failed to connect to etcd:", err)
		return err
	}
	defer client.Close()

	kv := clientv3.NewKV(client)
	key = fmt.Sprintf("/%s", key)
	// 读取键值
	resp, err := kv.Get(context.Background(), key)
	if err != nil {
		log.Fatal("Failed to get value from etcd:", err)
		return err
	}

	// 处理响应
	if len(resp.Kvs) > 0 {
		latestValue := resp.Kvs[0]
		err = jsonx.Unmarshal(latestValue.Value, c)
		if err != nil {
			log.Fatal("Failed to get value from etcd:", err)
			return err
		}
	} else {
		log.Fatal("No value found for key:", key)
		return err
	}
	fmt.Printf("Get config from etcd \n")

	gos.GoSafe(func() {
		c.EtcdChange(addr, key)
	})

	return nil
}

// EtcdChange 动态监听etcd变化，实时更新数据
func (c *Config) EtcdChange(addr []string, key string) {

	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("etcd failed, err:%v ", err)
		}
	}()

	f := func(value []byte) {
		log.Printf("etcd change %s", value)
		// 仅仅修改配置数据，如果之前初始化的内容需要重新初始化才生效。
		err := jsonx.Unmarshal(value, c)
		if err != nil {
			log.Fatalf("JsonUnmarshal etcd failed, err:%v ", err)
			return
		}
	}

	c.Config.EtcdChangeWatch(addr, key, f)
}
