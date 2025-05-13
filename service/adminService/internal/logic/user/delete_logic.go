package user

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

func (l *DeleteLogic) Delete(req *types.UserDeleteReq) (resp *types.UserDeleteResp, err error) {

	_, err = l.svcCtx.C.UserModel.DeleteByEx(l.ctx, goqu.Ex{"id": req.Id})

	if err != nil {
		logs.Err(l.ctx, "DeleteByEx: %v", err, logs.Flag("Delete"))
		return nil, err
	}

	resp = &types.UserDeleteResp{
		ApiBaseResp: types.ApiBaseResp{
			Code:    consts.Info_OK,
			Message: "success",
		},
	}
	return
}
