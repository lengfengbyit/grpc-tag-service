package main

import "context"

type Auth struct {
	AppKey    string
	AppSecret string
}

// GetRequestMetadata 获取请求所需要的元数据
func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

// RequireTransportSecurity 是否需要基于TLS认证进行安全传输
func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func NewAuth() *Auth {
	return &Auth{
		AppKey:    "go-blog",
		AppSecret: "go-blog",
	}
}
