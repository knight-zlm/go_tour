package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/knight-zlm/tag-service/proto"
	"github.com/knight-zlm/tag-service/server"
)

var (
	grpcPort string
	httpPort string
)

func init() {
	flag.StringVar(&grpcPort, "grpc_port", "8009", "grpc启动端口号")
	flag.StringVar(&httpPort, "http_port", "9009", "http启动端口号")
	flag.Parse()
}

//protoc --go_out=plugins=grpc:. ./proto/*.proto 编译proto
func main() {
	errsChan := make(chan error)
	go func() {
		err := RunHttpServer(httpPort)
		if err != nil {
			errsChan <- err
		}
	}()
	go func() {
		err := RunGrpcServer(grpcPort)
		if err != nil {
			errsChan <- err
		}
	}()

	select {
	case err := <-errsChan:
		log.Fatalf("Run Server err: %v", err)
	}
}

func RunHttpServer(port string) error {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})
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

	return s.Serve(lis)
}
