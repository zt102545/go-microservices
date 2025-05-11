db=$1
param=$2
# 字符串风格格式为：demoname
lower_name=$(echo $param | tr -d '_')
# 字符串风格格式为：DemoName
model_name=$(echo $param | awk -F '_' '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2));}1' | tr -d ' ')
# 字符串风格格式为：demoName
struct_name=$(echo $model_name | awk '{print tolower(substr($0,1,1)) substr($0,2)}')

help()
{
   cat <<- EOF

帮助文档：
  Desc: 后台CRUD代码生成工具
  程序选项如下：
    mysql：参数1:数据库表名,参数2:是否需要生成model的文件，默认不填不生成，1表示生成
    postgres：参数1:表名,参数2:是否需要生成model的文件，默认不填不生成，1表示生成,参数3:数据库,参数4:模式
  Example：
    ./generate_admin.sh mysql emails 1
    ./generate_admin.sh postgres tbl_mini_games 1 mini_games minigame
EOF
   exit 0
}

if [ "$param" == "" ]; then
  echo "参数为空"
  help
fi

if [ "$3" == "1" ]; then
  echo "生成model"
  if [ "$db" == "mysql" ]; then
    if ! [ -e ./dao/mysql/source/${param}.sql ]; then
        echo "数据源文件不存在"
        help
    fi
    # 生成mysql的crud
    ./generate_code.sh mysql $param
  elif [ "$db" == "postgres" ]; then
    # 生成mysql的crud
    ./generate_code.sh postgres $param $4 $5
  fi
fi
# model文件路径
model_file=./dao/${db}/model/${lower_name}model_gen.go

# 使用sed提取结构体并赋值给变量
#struct=$(sed -n "/^\s*"${model_name}" struct {/,/^\s*}$/p" $model_file)
struct=$(sed -n "/"${model_name}" struct {/,/^[[:space:]]*}$/p" $model_file)

if [ "$struct" == "" ]; then
  echo "struct为空"
  help
fi

# 结构体特殊类型替换
struct=$(echo "$struct" | sed 's/sql\.NullInt64/int64/g')
struct=$(echo "$struct" | sed 's/sql\.NullString/string/g')
struct=$(echo "$struct" | sed 's/sql\.NullTime/string/g')
struct=$(echo "$struct" | sed 's/time\.Time/string/g')
struct=$(echo "$struct" | sed 's/db/json/g')
struct=$(echo "$struct" | sed 's/"`/,optional"`/g')
struct=$(echo "$struct" | sed 's/struct//g')
struct=$(echo "$struct" | sed 's/^\t//g')
struct=$(echo "$struct" | sed -E 's/ +/ /g')

filename=./proto/source/admin/${param}.api

if [ -e $filename ]; then
    echo "api文件已生成"
    help
fi

cat > $filename <<- EOF
syntax = "v1"

import "base.api"

@server(
    prefix:     /admin/${lower_name} // 路由前缀
    group: ${lower_name} // 路由分组
    timeout:    5s // 超时配置
    middleware: AuthInterceptor // 路由添加中间件
    maxBytes:   1048576 // 请求体大小控制，单位为 byte,goctl 版本 >= 1.5.0 才支持
)
service admin {
    @doc(
        summary: "列表"
    )
    @handler List
    post /list (${model_name}ListReq) returns (${model_name}ListResp)

    @doc(
        summary: "删除"
    )
    @handler Delete
    post /delete (${model_name}DeleteReq) returns (${model_name}DeleteResp)

    @doc(
        summary: "新增"
    )
    @handler Insert
    post /insert (${model_name}InsertReq) returns (${model_name}InsertResp)

    @doc(
        summary: "修改"
    )
    @handler Update
    post /update (${model_name}UpdateReq) returns (${model_name}UpdateResp)
}

type ${struct}

type ${model_name}ListReq {
    ApiBaseReq
    Id int64 \`json:"id,optional"\`
    pageNumber int64 \`json:"page_number,optional"\`//页码
    pageSize int64 \`json:"page_size,optional"\`    //每页数量
}
type ${model_name}ListResp {
    ApiBaseResp
    Data interface{} \`json:"data"\`
}

type ${model_name}DeleteReq {
    ApiBaseReq
    Id int64 \`json:"id"\`
}
type ${model_name}DeleteResp {
    ApiBaseResp
}

type ${model_name}InsertReq {
    ApiBaseReq
    Data $model_name \`json:"data"\`
}
type ${model_name}InsertResp {
    ApiBaseResp
}

type ${model_name}UpdateReq {
    ApiBaseReq
    Data $model_name \`json:"data"\`
}
type ${model_name}UpdateResp {
    ApiBaseResp
}
EOF

base_path=./proto/source/admin/admin.api
if ! grep -Fxq "import \"${param}.api\"" "$base_path"; then
        insert_command=$(printf 's/^\/\/ shell/import \"%s.api\"\\n&/' "$param")
        sed -i '' -e "$insert_command" "$base_path"
fi

# admin后台api生成
./generate_code.sh api admin


# 删除代码生成
delete_file=./service/adminService/internal/logic/${lower_name}/delete_logic.go
# 获取行号，用来判断是否需要修改
line_number=$(sed -n '$=' ${delete_file})
if [ "$line_number" -lt 32 ]; then
  cat > "$delete_file" <<- EOF
package ${lower_name}

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"go-microservices/utils/consts"
	"go-microservices/utils/logs"

	"go-microservices/service/adminService/internal/svc"
	"go-microservices/service/adminService/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.${model_name}DeleteReq) (resp *types.${model_name}DeleteResp, err error) {

	_, err = l.svcCtx.C.${model_name}Model.DeleteByEx(l.ctx, goqu.Ex{"id": req.Id})

	if err != nil {
		logs.Err(l.ctx, "DeleteByEx: %v", err, logs.Flag("Delete"))
		return nil, err
	}

	resp = &types.${model_name}DeleteResp{
		ApiBaseResp: types.ApiBaseResp{
			Code:    consts.Info_OK,
			Message: "success",
		},
	}
	return
}
EOF
else
  echo "inaccurate number of codes, ignored generation"
fi
gofmt -w $delete_file

# 新增代码生成
insert_file=./service/adminService/internal/logic/${lower_name}/insert_logic.go
# 获取行号，用来判断是否需要修改
line_number=$(sed -n '$=' ${insert_file})
if [ "$line_number" -lt 32 ]; then
  cat > "$insert_file" <<- EOF
package ${lower_name}

import (
	"context"
	"go-microservices/dao/${db}/model"
	"go-microservices/utils/consts"
	"go-microservices/utils/logs"
	"go-microservices/utils/utils"

	"go-microservices/service/adminService/internal/svc"
	"go-microservices/service/adminService/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InsertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增
func NewInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InsertLogic {
	return &InsertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InsertLogic) Insert(req *types.${model_name}InsertReq) (resp *types.${model_name}InsertResp, err error) {

	dbModel := &model.${model_name}{}
	err = utils.CopyFields(&req.Data, dbModel)
	if err != nil {
		logs.Err(l.ctx, "Insert: %v", err, logs.Flag("Insert"))
		return nil, err
	}

	_, err = l.svcCtx.C.${model_name}Model.Insert(l.ctx, dbModel)
	if err != nil {
		logs.Err(l.ctx, "Insert: %v", err, logs.Flag("Insert"))
		return nil, err
	}

	resp = &types.${model_name}InsertResp{
		ApiBaseResp: types.ApiBaseResp{
			Code:    consts.Info_OK,
			Message: "success",
		},
	}
	return
}
EOF
else
  echo "inaccurate number of codes, ignored generation"
fi
gofmt -w $insert_file

# 列表代码生成
list_file=./service/adminService/internal/logic/${lower_name}/list_logic.go
# 获取行号，用来判断是否需要修改
line_number=$(sed -n '$=' ${list_file})
if [ "$line_number" -lt 32 ]; then
  cat > "$list_file" <<- EOF
package ${lower_name}

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"go-microservices/utils/consts"
	"go-microservices/utils/logs"
	"go-microservices/utils/utils"

	"go-microservices/service/adminService/internal/svc"
	"go-microservices/service/adminService/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 列表
func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.${model_name}ListReq) (resp *types.${model_name}ListResp, err error) {

	ex := goqu.Ex{}
	if req.Id > 0 {
		ex["id"] = req.Id
	}

	pageNumber := uint(1)
	pageSize := uint(consts.PAGESIZE)
	if req.PageNumber > 0 {
		pageNumber = uint(req.PageNumber)
	}
	if req.PageSize > 0 {
		pageSize = uint(req.PageSize)
	}
	orderBy := map[string]int{
		"id": 1,
	}

	list, err := l.svcCtx.C.${model_name}Model.FindList(l.ctx, ex, []uint{pageNumber, pageSize}, orderBy)
	if err != nil {
		logs.Err(l.ctx, "FindList: %v", err, logs.Flag("List"))
		return nil, err
	}

	data := make([]types.${model_name}, 0)
  	for _, v := range *list {
  		one := types.${model_name}{}
  		_ = utils.CopyFieldsBack(&v, &one)
  		data = append(data, one)
  }

	resp = &types.${model_name}ListResp{
		ApiBaseResp: types.ApiBaseResp{
			Code:    consts.Info_OK,
			Message: "success",
		},
		Data: data,
	}

	return
}
EOF
else
  echo "inaccurate number of codes, ignored generation"
fi
gofmt -w $list_file

# 更新代码生成
update_file=./service/adminService/internal/logic/${lower_name}/update_logic.go
# 获取行号，用来判断是否需要修改
line_number=$(sed -n '$=' ${update_file})
if [ "$line_number" -lt 32 ]; then
  cat > "$update_file" <<- EOF
package ${lower_name}

import (
	"context"
	"go-microservices/dao/${db}/model"
	"go-microservices/utils/consts"
	"go-microservices/utils/logs"
	"go-microservices/utils/utils"

	"go-microservices/service/adminService/internal/svc"
	"go-microservices/service/adminService/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.${model_name}UpdateReq) (resp *types.${model_name}UpdateResp, err error) {

	dbModel := &model.${model_name}{}
	err = utils.CopyFields(&req.Data, dbModel)
	if err != nil {
		logs.Err(l.ctx, "Update: %v", err, logs.Flag("Update"))
		return nil, err
	}

	err = l.svcCtx.C.${model_name}Model.Update(l.ctx, dbModel)
	if err != nil {
		logs.Err(l.ctx, "Update: %v", err, logs.Flag("Update"))
		return nil, err
	}

	resp = &types.${model_name}UpdateResp{
		ApiBaseResp: types.ApiBaseResp{
			Code:    consts.Info_OK,
			Message: "success",
		},
	}
	return
}
EOF
else
  echo "inaccurate number of codes, ignored generation"
fi
gofmt -w $update_file