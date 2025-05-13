package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/zeromicro/go-zero/core/jsonx"
	"go-microservices/utils/consts"
	"go-microservices/utils/logs"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
	"strings"
	"time"
)

type AuthInterceptorMiddleware struct {
}

type HttpResponseWriter struct {
	http.ResponseWriter
	code int
	res  *bytes.Buffer
}

type HttpLog struct {
	ContentType   string `json:"ContentType"`
	Method        string `json:"Method"`
	Url           string `json:"Url"`
	Header        string `json:"Header"`
	Req           string `json:"Req"`
	Res           string `json:"Res"`
	Code          int    `json:"Code"`
	RemainingTime int64  `json:"RemainingTime"`
}

func NewAuthInterceptorMiddleware() *AuthInterceptorMiddleware {
	return &AuthInterceptorMiddleware{}
}

func (m *AuthInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := r.FormValue("user_id")
		traceId := uuid.NewV4().String()
		md := metadata.New(map[string]string{
			consts.CTX_USERID:  userId,  // 用户id
			consts.CTX_TRACEID: traceId, // 链路id
		})

		ctx := metadata.NewOutgoingContext(r.Context(), md)
		ctx = context.WithValue(ctx, consts.CTX_USERID, userId)
		ctx = context.WithValue(ctx, consts.CTX_TRACEID, traceId)
		ctx = context.WithValue(ctx, consts.CTX_STARTTIME, time.Now())

		// 通过新上下文创建一个新的请求
		r = r.WithContext(ctx)
		// json方式需要复制一份body
		var body []byte
		if strings.Contains(r.Header.Get("Content-Type"), "json") {
			body, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(body))
		}

		rw := &HttpResponseWriter{
			ResponseWriter: w,
			code:           http.StatusOK,
			res:            new(bytes.Buffer),
		}
		next(rw, r)

		// 打印请求响应参数
		writeLog(r, rw, body)
	}
}

func (r *HttpResponseWriter) Write(bs []byte) (int, error) {
	r.res.Write(bs)
	return r.ResponseWriter.Write(bs)
}

func (r *HttpResponseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *HttpResponseWriter) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

func writeLog(r *http.Request, rw *HttpResponseWriter, body []byte) {

	reqMap := make(map[string]interface{})

	contentType := r.Header.Get("Content-Type")
	// url参数
	for k, v := range r.URL.Query() {
		reqMap[k] = v
	}

	// 解析并打印表单数据
	if strings.Contains(contentType, "form") {
		if err := r.ParseForm(); err == nil {
			for k, v := range r.Form {
				reqMap[k] = v
			}
		}
	}

	// 解析并打印JSON数据
	if strings.Contains(contentType, "json") && len(body) > 0 {
		bodyMap := make(map[string]interface{})
		if err := jsonx.Unmarshal(body, &bodyMap); err != nil {
			logs.Err(r.Context(), "Failed to unmarshal JSON data: %v", err)
		} else {
			for k, v := range bodyMap {
				reqMap[k] = v
			}
		}
	}

	req, err := jsonx.MarshalToString(reqMap)
	if err != nil {
		logs.Err(r.Context(), "http log json marshal error:", err)
	}

	headerJson, headerJsonErr := jsonx.MarshalToString(r.Header)

	if headerJsonErr != nil {
		logs.Err(r.Context(), "http log json marshal error:", headerJsonErr)
	}

	// 获取截止时间
	remainingTime := int64(0)
	deadline := r.Context().Value(consts.CTX_STARTTIME)
	if t, ok := deadline.(time.Time); ok {
		remainingTime = time.Now().Sub(t).Milliseconds()
	}

	httpLog := &HttpLog{
		ContentType:   contentType,
		Method:        r.Method,
		Url:           r.URL.Path,
		Header:        headerJson,
		Req:           req,
		Res:           rw.res.String(),
		Code:          rw.code,
		RemainingTime: remainingTime,
	}
	log, err := json.MarshalIndent(httpLog, "", "  ")
	if err != nil {
		logs.Err(r.Context(), "http log json marshal error:", err)
	}

	if rw.code >= 400 {
		logs.Err(r.Context(), strings.ReplaceAll(string(log), "\n", "\n\n"), logs.Flag(r.URL.Path))
	} else {
		logs.Info(r.Context(), strings.ReplaceAll(string(log), "\n", "\n\n"), logs.Flag(r.URL.Path))
	}
}
