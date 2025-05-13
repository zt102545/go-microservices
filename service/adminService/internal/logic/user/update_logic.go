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

func (l *UpdateLogic) Update(req *types.UserUpdateReq) (resp *types.UserUpdateResp, err error) {

	dbModel := &model.User{}
	err = utils.CopyFields(&req.Data, dbModel)
	if err != nil {
		logs.Err(l.ctx, "Update: %v", err, logs.Flag("Update"))
		return nil, err
	}

	err = l.svcCtx.C.UserModel.Update(l.ctx, dbModel)
	if err != nil {
		logs.Err(l.ctx, "Update: %v", err, logs.Flag("Update"))
		return nil, err
	}

	resp = &types.UserUpdateResp{
		ApiBaseResp: types.ApiBaseResp{
			Code:    consts.Info_OK,
			Message: "success",
		},
	}
	return
}
