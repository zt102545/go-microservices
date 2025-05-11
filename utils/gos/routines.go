package gos

import (
	"context"
	"fmt"
	"go-microservices/utils/consts"
	"go-microservices/utils/logs"
	"google.golang.org/grpc/metadata"
	"runtime/debug"
	"time"
)

func GoSafe(fn func()) {
	go runSafe(fn)
}

func CopyContextWithTimeout(parentCtx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {

	childCtx, childCancel := context.WithTimeout(context.Background(), timeout)

	if md, ok := metadata.FromIncomingContext(parentCtx); ok {
		childCtx = metadata.NewIncomingContext(childCtx, md)
	} else {
		childCtx = context.WithValue(childCtx, consts.CTX_USERID, parentCtx.Value(consts.CTX_USERID))
		childCtx = context.WithValue(childCtx, consts.CTX_TRACEID, parentCtx.Value(consts.CTX_TRACEID))
	}
	return childCtx, childCancel
}

func runSafe(fn func()) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("threading recover error is \n %+v\n %s", p, debug.Stack(), logs.Flag("GoSafe"))
		}
	}()
	fn()
}
