package errcode

import (
	pb "go-tour/grpc-tag-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TogRPCError 将业务错误转换为grpc错误
func TogRPCError(err *Error) error {
	// 将业务的错误信息放入grpc错误信息的details属性中
	pbErr := &pb.Error{Code: int32(err.Code()), Message: err.Msg()}
	s, _ := status.New(TogRPCCode(err.Code()), err.Msg()).WithDetails(pbErr)
	return s.Err()
}

// TogRPCCode 将业务错误码转换成grpc错误码
func TogRPCCode(code int) (statusCode codes.Code) {
	switch code {
	case Fail.Code():
		statusCode = codes.Internal
	case InvalidParams.Code():
		statusCode = codes.InvalidArgument
	case Unauthorized.Code():
		statusCode = codes.Unauthenticated
	case AccessDenied.Code():
		statusCode = codes.PermissionDenied
	case DeadlineExceeded.Code():
		statusCode = codes.DeadlineExceeded
	case NotFound.Code():
		statusCode = codes.NotFound
	case LimitExceed.Code():
		statusCode = codes.ResourceExhausted
	case MethodNotAllowed.Code():
		statusCode = codes.Unauthenticated
	default:
		statusCode = codes.Unknown
	}
	return
}

type Status struct {
	*status.Status
}

// FromError 获取原始的错误消息
func FromError(err error) *Status {
	s, _ := status.FromError(err)
	return &Status{s}
}

func TogRPCStatus(code int, msg string) *Status {
	pbErr := &pb.Error{Code: int32(code), Message: msg}
	s, _ := status.New(TogRPCCode(code), msg).WithDetails(pbErr)
	return &Status{s}
}
