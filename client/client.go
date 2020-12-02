package main

import (
	"context"
	"flag"
	"fmt"
	pb "go-tour/grpc-tag-service/proto"
	"google.golang.org/grpc"
	"log"
)

var host string
var port string

func init() {
	flag.StringVar(&host, "host", "localhost", "Server listen ip")
	flag.StringVar(&port, "port", "8001", "Server listen port")
	flag.Parse()
}

func getServerAddr() string {
	return fmt.Sprintf("%s:%s", host, port)
}

// grpc 客户端
func main() {
	ctx := context.Background()
	clientConn, _ := GetClientConn(ctx, getServerAddr(), nil)
	defer clientConn.Close()

	tagServerClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServerClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("Client GetTagList err: %v", err)
	}
	log.Printf("Client GetTagList resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}
