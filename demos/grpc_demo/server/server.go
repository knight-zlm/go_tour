package main

import (
	"context"
	"flag"
	"net"

	pb "grpc-demo/proto"

	"google.golang.org/grpc"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8008", "启动端口号")
	flag.Parse()
}

type GreeterServer struct {
}

func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello.word"}, nil
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis, _ := net.Listen("tcp", ":"+port)
	server.Serve(lis)
}
