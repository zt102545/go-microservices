package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-microservices/service/adminService/internal/logic/user"
	"go-microservices/service/adminService/internal/svc"
	"go-microservices/service/adminService/internal/types"
)

// 删除
func DeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewDeleteLogic(r.Context(), svcCtx)
		resp, err := l.Delete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
