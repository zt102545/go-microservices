db=$1
param=$2
# 字符串风格格式为：DemoName
model_name=$(echo "${param}" | awk -F '_' '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2));}1' | tr -d ' ')
# 字符串风格格式为：demoName
struct_name=$(echo "${model_name}" | awk '{print tolower(substr($0,1,1)) substr($0,2)}')
help()
{
   cat <<- EOF

帮助文档：
  Desc: 基于goctl指令生成代码
  程序选项如下：
    api：生成api框架和swagger文档,参数1:api文件名(不带扩展名)
    rpc：生成rpc框架和proto文件，参数1:rpc文件名(不带扩展名)
    proto 只生成base目录下的proto文件，注意路径，参数1:base
    mysql：生成model和crud代码,参数1:源文件名(不带扩展名)
    postgres ：生成model和crud代码,参数1:表名,参数2:数据库,参数3:模式
    docker：打包docker镜像,参数1:main文件名(不带扩展名),参数2:版本号,参数3:镜像推送环境(dev/stag)
  Example:
    ./generate_code.sh api gateway
    ./generate_code.sh rpc user
    ./generate_code.sh proto base
    ./generate_code.sh mysql users
    ./generate_code.sh postgres tbl_mini_games mini_games minigame
    ./generate_code.sh docker users 1.0.0 dev
EOF
   exit 0
}

# api框架代码生成
gen_api()
{
    # 获取要修改的main文件名
    main=./service/${param}Service/${param}.go
    # 获取行号，用来判断是否需要修改
    line_number=$(sed -n '$=' ${main})
    if [ "$line_number" -lt 40 ]; then
      cat > "$main" <<- EOF
package main

import (
    "flag"
    "fmt"
    _ "github.com/apache/skywalking-go"
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/rest"
    "go-microservices/service/${param}Service/internal/config"
    "go-microservices/service/${param}Service/internal/handler"
    "go-microservices/service/${param}Service/internal/svc"
    "os"
    "strings"
)

var configFile = flag.String("f", "./service/${param}Service/etc/${param}.yaml", "the config file")

func main() {

    flag.Parse()
    var c config.Config
    if len(os.Args) > 1 {
      etcd := strings.Split(os.Args[1], ",")
      err := c.NewConfigFromEtcd(etcd, "${param}Config")
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
EOF
    else
      echo "inaccurate number of codes, ignored generation"
    fi

    # 获取要修改的config文件名
    config=./service/${param}Service/internal/config/config.go
    # 获取行号，用来判断是否需要修改
    line_number=$(sed -n '$=' ${config})
    if [ "$line_number" -lt 10 ]; then
      cat > "$config" <<- EOF
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
EOF
    else
      echo "inaccurate number of codes, ignored generation"
    fi

    # 获取要修改的service文件名
    service=./service/${param}Service/internal/svc/service_context.go
    # 获取行号，用来判断是否需要修改
    line_number=$(sed -n '$=' ${service})
    if [ "$line_number" -lt 20 ]; then
      cat > "$service" <<- EOF
package svc

import (
    "github.com/zeromicro/go-zero/rest"
    "go-microservices/com"
    "go-microservices/service/${param}Service/internal/config"
    "go-microservices/service/${param}Service/internal/middleware"
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
EOF
    else
      echo "inaccurate number of codes, ignored generation"
    fi
}

# rpc框架代码生成
gen_rpc()
{
    # 获取要修改的main文件名
    main=./service/${param}Service/${param}.go
    # 获取行号，用来判断是否需要修改
    line_number=$(sed -n '$=' ${main})
    if [ "$line_number" -lt 40 ]; then
      cat > "$main" <<- EOF
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

var configFile = flag.String("f", "./service/${param}Service/etc/${param}.yaml", "the config file")

func main() {

    flag.Parse()
    var c config.Config
    if len(os.Args) > 1 {
      etcd := strings.Split(os.Args[1], ",")
      err := c.NewConfigFromEtcd(etcd, "${param}Config")
      if err != nil {
        panic(err)
      }
    } else {
      conf.MustLoad(*configFile, &c)
    }

    ctx := svc.NewServiceContext(c)
    serverConf := c.GetRpcServerConf()
    s := zrpc.MustNewServer(serverConf, func(grpcServer *grpc.Server) {
      ${param}.Register${model_name}Server(grpcServer, server.New${model_name}Server(ctx))

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
EOF
    else
      echo "inaccurate number of codes, ignored generation"
    fi

    # 获取要修改的main文件名
    config=./service/${param}Service/internal/config/config.go
    # 获取行号，用来判断是否需要修改
    line_number=$(sed -n '$=' ${config})
    if [ "$line_number" -lt 10 ]; then
      cat > "$config" <<- EOF
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
EOF
    else
      echo "inaccurate number of codes, ignored generation"
    fi

    # 获取要修改的service文件名
    service=./service/${param}Service/internal/svc/service_context.go
    # 获取行号，用来判断是否需要修改
    line_number=$(sed -n '$=' ${service})
    if [ "$line_number" -lt 20 ]; then
      cat > "$service" <<- EOF
package svc

import (
    "github.com/zeromicro/go-zero/rest"
    "go-microservices/com"
    "go-microservices/service/${param}Service/internal/config"
    "go-microservices/service/${param}Service/internal/middleware"
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
EOF
    else
      echo "inaccurate number of codes, ignored generation"
    fi
}

# 生成CRUD的方法
gen_crud()
{
    # 获取要修改的文件名
    filename=./dao/"${db}"/model/$(echo "${param}" | tr -d '_')model.go
    # 获取行号，用来判断是否需要修改
    line_number=$(sed -n '$=' ${filename})
    if [ "$line_number" -lt 30 ]; then
      cat > "$filename" <<- EOF
package model

import (
    "context"
    "database/sql"
    "errors"
    "github.com/doug-martin/goqu/v9"
    _ "github.com/doug-martin/goqu/v9/dialect/${db}"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "strings"
    "go-microservices/utils/utils"
)

var _ ${model_name}Model = (*custom${model_name}Model)(nil)

type (
    // ${model_name}Model is an interface to be customized, add more methods here,
    // and implement the added methods in custom${model_name}Model.
    ${model_name}Model interface {
        ${struct_name}Model
        withSession(session sqlx.Session) ${model_name}Model
        // GetTableName 获取表名
        GetTableName() string
        // GetCount 根据条件获取数量
        GetCount(ctx context.Context, ex goqu.Expression) (int64, error)
        // FindList 根据条件获取列表，排序：map[string]int{"字段":0/1(0-升序(ASC)；1-降序(DESC))}；分页：[]uint{页码，每页条数}
        FindList(ctx context.Context, ex goqu.Expression, optionalParams ...any) (*[]${model_name}, error)
        // FindOnly 根据条件获取单条数据，0-升序(ASC)；1-降序(DESC)
        FindOnly(ctx context.Context, ex goqu.Expression, order ...map[string]int) (*${model_name}, error)
        // InsertOnly 插入单条数据
        InsertOnly(ctx context.Context, row *${model_name}, tx ...*sql.Tx) (sql.Result, error)
        // BatchInsert 批量插入
        BatchInsert(ctx context.Context, rows []*${model_name}, tx ...*sql.Tx) (sql.Result, error)
        // UpdateByEx 根据条件更新
        UpdateByEx(ctx context.Context, record goqu.Record, ex goqu.Expression, tx ...*sql.Tx) (sql.Result, error)
        // DeleteByEx 根据条件删除数据
        DeleteByEx(ctx context.Context, ex goqu.Expression, tx ...*sql.Tx) (sql.Result, error)
    }

    custom${model_name}Model struct {
        *default${model_name}Model
    }
)

// New${model_name}Model returns a model for the database table.
func New${model_name}Model(conn sqlx.SqlConn) ${model_name}Model {
    return &custom${model_name}Model{
      default${model_name}Model: new${model_name}Model(conn),
    }
}

func (m *custom${model_name}Model) withSession(session sqlx.Session) ${model_name}Model {
    return New${model_name}Model(sqlx.NewSqlConnFromSession(session))
}

// GetTableName 获取表名
func (m *custom${model_name}Model) GetTableName() string {
    return utils.SetTable(m.table)
}

// GetCount 根据条件获取数量
func (m *custom${model_name}Model) GetCount(ctx context.Context, ex goqu.Expression) (int64, error) {
    query, _, err := goqu.Dialect("${db}").Select(goqu.COUNT(1)).From(utils.SetTable(m.table)).Where(ex).ToSQL()
    if err != nil {
        return 0, err
    }
    var resp int64
    err = m.conn.QueryRowCtx(ctx, &resp, query)
    if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
        return 0, err
    }
    return resp, nil
}

// FindList 根据条件获取列表，排序：map[string]int{"字段":0/1(0-升序(ASC)；1-降序(DESC))}；分页：[]uint{页码，每页条数}
func (m *custom${model_name}Model) FindList(ctx context.Context, ex goqu.Expression, optionalParams ...any) (*[]${model_name}, error) {
    sql := goqu.Dialect("${db}").Select(&${model_name}{}).From(utils.SetTable(m.table)).Where(ex)
    if len(optionalParams) > 0 {
        for _, param := range optionalParams {
            // 排序
            if v, ok := param.(map[string]int); ok {
                for key, value := range v {
                    if value > 0 {
                        sql = sql.OrderAppend(goqu.C(key).Desc())
                    } else {
                        sql = sql.OrderAppend(goqu.C(key).Asc())
                    }
                }
            }
            // 分页
            if v, ok := param.([]uint); ok {
                if len(v) == 2 {
                    sql = sql.Offset((v[0] - 1) * v[1]).Limit(v[1])
                }
            }
        }
    }
    query, _, err := sql.ToSQL()
    query = strings.ReplaceAll(query, \`""\`, \`"\`)
    if err != nil {
        return nil, err
    }
    var resp []${model_name}
    err = m.conn.QueryRowsCtx(ctx, &resp, query)
    if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
        return nil, err
    }
    return &resp, nil
}

// FindOnly 根据条件获取单条数据，0-升序(ASC)；1-降序(DESC)
func (m *custom${model_name}Model) FindOnly(ctx context.Context, ex goqu.Expression, order ...map[string]int) (*${model_name}, error) {
    sql := goqu.Dialect("${db}").Select(&${model_name}{}).From(utils.SetTable(m.table)).Where(ex)
    if len(order) > 0 {
        for key, value := range order[0] {
          if value > 0 {
              sql = sql.OrderAppend(goqu.C(key).Desc())
          } else {
              sql = sql.OrderAppend(goqu.C(key).Asc())
          }
        }
    }
    query, _, err := sql.Limit(1).ToSQL()
    query = strings.ReplaceAll(query, \`""\`, \`"\`)
    if err != nil {
        return nil, err
    }
    var resp ${model_name}
    err = m.conn.QueryRowCtx(ctx, &resp, query)
    switch err {
    case nil:
        return &resp, nil
    case sqlx.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

// InsertOnly 插入单条数据
func (m *custom${model_name}Model) InsertOnly(ctx context.Context, row *${model_name}, tx ...*sql.Tx) (sql.Result, error) {
    query, _, err := goqu.Dialect("${db}").Insert(utils.SetTable(m.table)).Rows(row).ToSQL()
    if err != nil {
        return nil, err
    }
    var result sql.Result
    if len(tx) > 0 {
        result, err = tx[0].ExecContext(ctx, query)
    } else {
        result, err = m.conn.ExecCtx(ctx, query)
    }
    return result, err
}

// BatchInsert 批量插入
func (m *custom${model_name}Model) BatchInsert(ctx context.Context, rows []*${model_name}, tx ...*sql.Tx) (sql.Result, error) {
    query, _, err := goqu.Dialect("${db}").Insert(utils.SetTable(m.table)).Rows(rows).ToSQL()
    if err != nil {
        return nil, err
    }
    var result sql.Result
    if len(tx) > 0 {
        result, err = tx[0].ExecContext(ctx, query)
    } else {
        result, err = m.conn.ExecCtx(ctx, query)
    }
    return result, err
}

// UpdateByEx 根据条件更新
func (m *custom${model_name}Model) UpdateByEx(ctx context.Context, record goqu.Record, ex goqu.Expression, tx ...*sql.Tx) (sql.Result, error) {
    query, _, err := goqu.Dialect("${db}").Update(utils.SetTable(m.table)).Set(record).Where(ex).ToSQL()
    if err != nil {
        return nil, err
    }
    var result sql.Result
    if len(tx) > 0 {
        result, err = tx[0].ExecContext(ctx, query)
    } else {
        result, err = m.conn.ExecCtx(ctx, query)
    }
    return result, err
}

// DeleteByEx 根据条件删除数据
func (m *custom${model_name}Model) DeleteByEx(ctx context.Context, ex goqu.Expression, tx ...*sql.Tx) (sql.Result, error) {
    query, _, err := goqu.Dialect("${db}").Delete(utils.SetTable(m.table)).Where(ex).ToSQL()
    if err != nil {
        return nil, err
    }
    var result sql.Result
    if len(tx) > 0 {
        result, err = tx[0].ExecContext(ctx, query)
    } else {
        result, err = m.conn.ExecCtx(ctx, query)
    }
    return result, err
}
EOF
    else
      echo "inaccurate number of codes, ignored generation"
    fi
}

dockerfile()
{
   cat > "Dockerfile" <<- EOF
FROM golang:1.22-alpine AS builder

LABEL stage=gobuilder

ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

# 设置工作目录
WORKDIR /build

# 加载依赖
ADD go.mod .
ADD go.sum .
RUN go mod download

# 复制源代码
COPY . .

# 静态编译Go程序
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -toolexec="/build/doc/tool/skywalking-go-agent" -a -o app ./service/${struct_name}Service/${param}.go

# 第二阶段：运行时镜像，使用空镜像scratch或者alpine
FROM --platform=linux/amd64 alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# 设置工作目录
WORKDIR /app

# 复制编译好的二进制文件到运行时镜像
COPY --from=builder /build/app .
COPY --from=builder /build/i18n /app/i18n
COPY --from=builder /build/key/google_play_service_account_key.json /etc/ssl/certs/google_play_service_account_key.json
COPY --from=builder /build/key/SubscriptionKey_23HSTRGFC3.p8 /etc/ssl/certs/SubscriptionKey_23HSTRGFC3.p8
COPY --from=builder /build/migrations /app/migrations

# 设置 SkyWalking Agent 的配置
ENV SW_AGENT_NAME=go-${param}
ENV SW_AGENT_REPORTER_GRPC_BACKEND_SERVICE=skywalking-oap.skywalking.svc.cluster.local:11800
ENV SW_AGENT_REPORTER_GRPC_CDS_FETCH_INTERVAL=-1
# 屏蔽日志上传
ENV SW_AGENT_LOG_TRACING_ENABLE=false
ENV SW_AGENT_LOG_REPORTER_ENABLE=false
# 屏蔽etcd心跳,测试了并没有效果
ENV SW_AGENT_TRACE_IGNORE_PATH=etcdserverpb**,/etcdserverpb/**,**etcdserverpb**,**/etcdserverpb/**

# 运行程序
ENTRYPOINT ["./app"]
# 设置CMD指令来指定参数，默认测试环境的etcd
CMD ["16.162.220.93:2379"]
EOF
}

if [ "$1" == "api" ]; then
    # api网关生成
    goctl api go -api ./proto/source/${param}/${param}.api -dir ./service/${param}Service -style go_zero
    # swagger文档生成
    goctl api plugin -plugin goctl-swagger="swagger -filename ${param}.json -host 127.0.0.1:8888" -api ./proto/source/${param}/${param}.api -dir ./doc/swagger/etc
    gen_api
    echo "Done."
elif [ "$1" == "rpc" ]; then
    goctl rpc protoc ./proto/source/${param}.proto --go_out=./proto/generate --go-grpc_out=./proto/generate --zrpc_out=./service/${struct_name}Service --style go_zero
    # 替换omitempty，避免json序列化忽略字段
    sed -i '' -e '/omitempty/s/,omitempty//g' ./proto/generate/${struct_name}/*.pb.go

    # 修改客户端文件包名
    path="./service/${struct_name}Service/${param}_client/"
    sed -i '' -e "s/package ${param}_client/package client/g" ${path}${param}.go
    # 将客户端文件移到./proto/client下，删除原来目录
    mv -f ${path}${param}.go ./proto/client
    rm -r ${path}
    gen_rpc
    echo "Done."
elif [ "$1" == "proto" ]; then
    # proto文件生成，base目录下
    protoc --go_out=.. --go-grpc_out=..  ./proto/source/${param}.proto
    sed -i '' -e '/omitempty/s/,omitempty//g' ./proto/generate/${param}/*.pb.go
    echo "Done."
elif [ "$1" == "mysql" ]; then
    # mysql生成代码
    goctl model mysql ddl --src ./dao/mysql/source/${param}.sql --dir ./dao/mysql/model -i '' --strict
    gen_crud
    echo "Done."
elif [ "$1" == "postgres" ]; then
    dbName=$3
    schema=$4
    # postgres生成代码
    goctl model pg datasource --url="postgres://postgres:yjuVnq2yW3XQgFitAZNf@43.198.115.11:30005/${dbName}?sslmode=disable" --schema="${schema}" --table="${param}" --dir ./dao/postgres/model -i '' --strict
    gen_crud
    echo "Done."
elif [ "$1" == "docker" ]; then
    version=$3
    env=$4
    # 推送地址，镜像名，根据需要修改
    tagName="io/${env}/go-microservices/${param}:${version}"
    dockerfile
    docker build --platform="linux/amd64" --no-cache -t "${param}" .
    docker tag "${param}" "${tagName}"
    docker push "${tagName}"
    docker rmi "${param}"
    docker rmi "${tagName}"
    rm -f Dockerfile
    echo "Done."
else
    echo "参数无效"
    help
fi