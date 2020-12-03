package auth

import (
	"context"
	"go-tour/grpc-tag-service/pkg/errcode"
	"google.golang.org/grpc/metadata"
)

type Auth struct {
}

func (a *Auth) GetAppKey() string {
	return "go-blog"
}

func (a *Auth) GetAppSecret() string {
	return "go-blog"
}

// Check 权限检查
func (a *Auth) Check(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)

	var appKey, appSecret string
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}

	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}

	return nil
}
