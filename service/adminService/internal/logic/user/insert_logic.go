package user

import (
	"context"
	"go-microservices/dao/mysql/model"
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

func (l *InsertLogic) Insert(req *types.UserInsertReq) (resp *types.UserInsertResp, err error) {

	dbModel := &model.User{}
	err = utils.CopyFields(&req.Data, dbModel)
	if err != nil {
		logs.Err(l.ctx, "Insert: %v", err, logs.Flag("Insert"))
		return nil, err
	}

	_, err = l.svcCtx.C.UserModel.Insert(l.ctx, dbModel)
	if err != nil {
		logs.Err(l.ctx, "Insert: %v", err, logs.Flag("Insert"))
		return nil, err
	}

	resp = &types.UserInsertResp{
		ApiBaseResp: types.ApiBaseResp{
			Code:    consts.Info_OK,
			Message: "success",
		},
	}
	return
}
