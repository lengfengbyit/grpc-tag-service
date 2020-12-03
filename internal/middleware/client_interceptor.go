package middleware

import (
	"context"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		// 设置默认超时时间，该超时时间是针对整条调用链路的
		defaultTimeout := 10 * time.Second
		ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
	}
	return ctx, cancel
}

// UnaryContextTimeout 设置一元RPC调用的超时时间
func UnaryContextTimeout() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		ctx, cancel := defaultContextTimeout(ctx)
		if cancel != nil {
			defer cancel()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// UnaryRetry 根据错误吗,进行接口重试请求
func UnaryRetry() grpc.UnaryClientInterceptor {
	return grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithMax(2), // 最大重试次数
		grpc_retry.WithCodes( // 设置对那些错误码进行重试操作
			codes.Unknown,
			codes.Internal,
			codes.DeadlineExceeded,
		),
	)
}

// StreamContextTimeout 设置流式RPC调用的超时时间
func StreamContextTimeout() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx, cancel := defaultContextTimeout(ctx)
		if cancel != nil {
			defer cancel()
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

// StreamRetry 设置流式RPC调用的重试机制
func StreamRetry() grpc.StreamClientInterceptor {
	return grpc_retry.StreamClientInterceptor(
		grpc_retry.WithMax(2),
		grpc_retry.WithCodes(
			codes.Unknown,
			codes.Internal,
			codes.DeadlineExceeded,
		),
	)
}
