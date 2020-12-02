package main

import (
	"flag"
	"fmt"
	"github.com/soheilhy/cmux"
	pb "go-tour/grpc-tag-service/proto"
	"go-tour/grpc-tag-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

// HTTP and GRPC
var port string

func init() {
	flag.StringVar(&port, "port", "8001", "server listen port")
	flag.Parse()
}

func main() {
	lis, err := RunTCPServer(port)
	if err != nil {
		log.Fatalf("Run tcp server err: %v", err)
	}

	m := cmux.New(lis)
	grpcLis := m.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldPrefixSendSettings(
			"content-type",
			"application/grpc",
		),
	)
	httpLis := m.Match(cmux.HTTP1Fast())

	grpcSer := RunGrpcServer()
	httpSer := RunHttpServer()
	go grpcSer.Serve(grpcLis)
	go httpSer.Serve(httpLis)

	err = m.Serve()
	if err != nil {
		log.Fatalf("Run serve err: ", err)
	}
}

func RunTCPServer(port string) (net.Listener, error) {
	fmt.Println("Server (http and grpc) listen http://localhost:" + port)
	return net.Listen("tcp", ":"+port)
}

func RunGrpcServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	return s
}

func RunHttpServer() *http.Server {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})
	return &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
}
