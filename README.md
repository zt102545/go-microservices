# go-microservices

> 基于go-zero的微服务框架

## 技术栈
    框架文档：https://go-zero.dev/docs/tasks

+ 后端框架：go-zero;grpc;sqlx;go-redis;
+ 服务注册与发现：etcd;
+ 系统配置：etcd;
+ 数据库：mysql;
+ 消息队列：kafka;
+ 缓存：redis;cache;
+ 日志：elk;
+ 监控：skywalking;
+ 部署：docker;k8s;

## 安装

1. 搭建 [Golang](https://golang.google.cn/) 环境 ,go版本1.22.2
2. 安装 goctl，版本1.6.5
  ```shell
go install github.com/zeromicro/go-zero/tools/goctl@v1.6.5 
# 查看go环境变量地址，该目录下如果没有goctl，手动拷贝到该目录
go env GOPATH

# 添加环境变量文件
touch ~/.zshrc
# 添加以下内容
export GOPATH=$HOME/go
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
# 生效环境变量
source ~/.zshrc
```

3. 安装rpc
```shell
goctl env check --install --verbose --force 
# 查看go环境变量地址，该目录下如果没有goctl，手动拷贝到该目录
go env GOPATH
```

4. 生成go.mod，已有可以忽略
```shell
go mod init go-microservices
```

5. 搭建ETCD
```shell
brew install etcd
# 设置开机启动
brew services start etcd

# docker方式
docker pull bitnami/etcd:latest
docker run -d --name etcd -p 2379:2379 -p 2380:2380 -e ALLOW_NONE_AUTHENTICATION=yes bitnami/etcd
```
[ETCD客户端下载](https://github.com/workpieces/etcdwp/releases)

6. 安装kafka
```shell
brew install kafka
# 设置开机启动
brew services start zookeeper
brew services start kafka
```
[kafka客户端下载](https://www.kafkatool.com/download.html)
7. swagger
```shell
go install github.com/zeromicro/goctl-swagger@latest
# 查看go环境变量地址，该目录下如果没有goctl，手动拷贝到该目录
go env GOPATH
```

## 操作指令

+ 创建API/RPC服务示例

```shell
# 基于.proto/source下的文件创建api/rpc服务
# 生成api后台框架
./generate_code.sh api admin
# 生成api网关
./generate_code.sh api gateway
# 生成rpc服务
./generate_code.sh rpc demo
```

+ 生成服务代码示例

```shell
# 基于dao/mysql/source下的sql文件生成后台crud基本代码
./generate_admin.sh mysql user 1

# 基于.proto/source下的文件生成代码
# 仅生成rpc文件，不生成框架
./generate_code.sh proto base
# 生成api网关
./generate_code.sh api gateway
# 生成rpc服务
./generate_code.sh rpc demo
```

+ 启动服务

```shell
go run 服务名称.go -f 配置文件地址
# 1. 启动 网管服务
go run user.go -f ./service/gateway/etc/gateway-api.yaml

# 或者直接配置程序实参后直接启动
-f ./service/gatewayService/etc/gatewayservice-api.yaml
```

# 目录结构

- go-microservices
  - com【业务公共代码】 
  - dao【model层】
    - mysql
      - model【生成curd和model的目录】
      - source【表结构】
    - redis
      - key【redis key】
      - model【redis操作逻辑】
  - doc【文档】
    - swagger【接口文档】
  - proto
    - client【rpc客户端连接】
    - generate【生成rpc目录】
    - source【proto源文件】
  - service【微服务集合】
    - demoService【服务名称】
      - etc【配置文件】
      - internal
        - config【配置结构体】
        - logic【业务逻辑】
        - server【service】
        - svc【初始化依赖】
      - demo.go【main入口】
  - utils【公共代码模块】
    - config【配置初始化】
    - consts【常量】
    - gos【协程】
    - kafka【kafka相关】
    - logs【日志相关】
    - mysql【mysql】
    - redis【redis】
    - utils【公共代码】