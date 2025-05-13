package user

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

func (l *ListLogic) List(req *types.UserListReq) (resp *types.UserListResp, err error) {

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

	list, err := l.svcCtx.C.UserModel.FindList(l.ctx, ex, []uint{pageNumber, pageSize}, orderBy)
	if err != nil {
		logs.Err(l.ctx, "FindList: %v", err, logs.Flag("List"))
		return nil, err
	}

	data := make([]types.User, 0)
	for _, v := range *list {
		one := types.User{}
		_ = utils.CopyFieldsBack(&v, &one)
		data = append(data, one)
	}

	resp = &types.UserListResp{
		ApiBaseResp: types.ApiBaseResp{
			Code:    consts.Info_OK,
			Message: "success",
		},
		Data: data,
	}

	return
}
