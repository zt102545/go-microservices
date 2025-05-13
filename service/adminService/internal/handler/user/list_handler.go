package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-microservices/service/adminService/internal/logic/user"
	"go-microservices/service/adminService/internal/svc"
	"go-microservices/service/adminService/internal/types"
)

// 列表
func ListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewListLogic(r.Context(), svcCtx)
		resp, err := l.List(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
