package main

import (
	"flag"
	"fmt"
	pb "go-tour/grpc-tag-service/proto"
	"go-tour/grpc-tag-service/server"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc/reflection"
)

var grpcPort string
var httpPort string

func init() {
	flag.StringVar(&grpcPort, "grpc_port", "8001", "server listen grpc port")
	flag.StringVar(&httpPort, "http_port", "9001", "server listen http port")
	flag.Parse()
}

func main() {
	errs := make(chan error)
	go func() {
		err := RunHttpServer(httpPort)
		if err != nil {
			errs <- err
		}
	}()

	go func() {
		err := RunGrpcServer(grpcPort)
		if err != nil {
			errs <- err
		}
	}()

	select {
	case err := <-errs:
		log.Fatalf("Run server err: %s", err)
	}
}

func RunHttpServer(port string) error {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

	fmt.Println("Http server listen: http://localhost:" + port)
	return http.ListenAndServe(":"+port, serveMux)
}

func RunGrpcServer(port string) error {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	fmt.Println("GRPC Server listen: http://localhost:" + port)
	return s.Serve(lis)
}
