package demo

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-microservices/service/gatewayService/internal/logic/demo"
	"go-microservices/service/gatewayService/internal/svc"
	"go-microservices/service/gatewayService/internal/types"
)

// ping
func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PingReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := demo.NewPingLogic(r.Context(), svcCtx)
		resp, err := l.Ping(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
