package middleware

import (
	"context"
	"go-tour/grpc-tag-service/pkg/errcode"
	"google.golang.org/grpc"
	"log"
	"runtime/debug"
	"time"
)

const dateLayout = "2006-01-02 15:04:05"

// AccessLog 访问日志拦截器
func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	reqLog := "[access request log]: method: %s, begin time: %s, request: %v"
	beginTime := time.Now().Local().Format(dateLayout)
	log.Printf(reqLog, info.FullMethod, beginTime, req)

	resp, err := handler(ctx, req)

	responseLog := "[access response log]: method: %s, begin time: %s, end time: %d, response: %v"
	endTime := time.Now().Local().Format(dateLayout)
	log.Printf(responseLog, info.FullMethod, beginTime, endTime, resp)

	return resp, err
}

// ErrorLog 错误日志记录
func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	resp, err := handler(ctx, req)
	if err != nil {
		errLog := "[error log]: method: %s, code: %v, message: %v, details: %v"
		s := errcode.FromError(err)
		log.Printf(errLog, info.FullMethod, s.Code(), s.Err().Error(), s.Details())
	}

	return resp, err
}

func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if e := recover(); e != nil {
			recoverLog := "[recovery log]: method: %s, message: %v, stack: %s"
			log.Printf(recoverLog, info.FullMethod, e, string(debug.Stack()[:]))
		}
	}()

	return handler(ctx, req)
}
