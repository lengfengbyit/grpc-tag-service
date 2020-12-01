package main

import (
	"flag"
	"fmt"
	pb "go-tour/grpc-tag-service/proto"
	"go-tour/grpc-tag-service/server"
	"google.golang.org/grpc"
	"log"
	"net"

	"google.golang.org/grpc/reflection"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8080", "server listen port")
	flag.Parse()
}

func main() {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	// 注册服务，方便使用grpcurl调试
	// grpcurl -plaintext localhost:8080 list 查看服务列表
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen err: %s", err)
	}

	fmt.Println("Server listen: http://127.0.0.1:" + port)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("server.Serve err: %s", err)
	}
}
