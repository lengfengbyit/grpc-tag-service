package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go-tour/grpc-tag-service/internal/middleware"
	"go-tour/grpc-tag-service/pkg/swagger"
	"go-tour/grpc-tag-service/pkg/tracer"
	pb "go-tour/grpc-tag-service/proto"
	"go-tour/grpc-tag-service/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"path"
	"strings"
)

// HTTP and GRPC
var port string

type httpError struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func init() {
	flag.StringVar(&port, "port", "8001", "server listen port")
	flag.Parse()

	// 设置 tracer 配置
	err := setupTracer()
	if err != nil {
		log.Fatalf("setupTracer err: %v",err)
	}
}

func main() {

	fmt.Println("Server (grpc and http) listen: http://localhost:" + port)
	err := RunServer(port)
	if err != nil {
		log.Fatalf("Run serve err: ", err)
	}
}

func RunServer(port string) error {
	httpMux := runHttpServer()
	grpcSer := runGrpcServer()
	gatewayMux := runGrpcGatewayServer()

	httpMux.Handle("/", gatewayMux)
	return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcSer, httpMux))
}

func runGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{
		// 注册中间件
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middleware.Recovery,
			middleware.ServerTracing,
			middleware.AccessLog,
			middleware.ErrorLog,
		)),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	return s
}

func runHttpServer() *http.ServeMux {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

	// 访问 swagger-ui
	prefix := "/swagger-ui/"
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	serverMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	serverMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("url: %s", r.URL.Path)
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
		p = path.Join("proto", p)
		http.ServeFile(w, r, p)
	})
	return serverMux
}

func runGrpcGatewayServer() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + port
	// 使用自定义的错误处理方法
	runtime.HTTPError = grpcGatewayError
	gwmux := runtime.NewServeMux()
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dopts)
	return gwmux
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	// h2c 允许通过明文TCP运行HTTP/2协议
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.ProtoMajor 客户端请求的协议版本号 http/1.1 or http/2
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

// grpcGatewayError 将 grpc 错误转换为 json 格式的错误，返回给客户端
func grpcGatewayError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler,
	w http.ResponseWriter, _ *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	httpError := httpError{Code: int32(s.Code()), Message: s.Message()}
	details := s.Details()
	for _, detail := range details {
		if v, ok := detail.(*pb.Error); ok {
			httpError.Code = v.Code
			httpError.Message = v.Message
		}
	}

	resp, _ := json.Marshal(httpError)
	w.Header().Set("Content-Type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code())) // 将grpc错误码转换成http的状态码
	_, _ = w.Write(resp)
}

func setupTracer() error {
	_, _, err := tracer.NewJaegerTracer(
		"grpc-tag-service",
		"127.0.0.1:6831",
	)
	if err != nil {
		return err
	}

	return nil
}
