package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
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

func (s *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&pb.HelloReply{Message: fmt.Sprintf("hello.%d", n)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GreeterServer) SayRecord(stream pb.Greeter_SayRecordServer) error {
	for true {
		recv, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{Message: "say.record"})
		}
		if err != nil {
			return err
		}

		log.Printf("recv: %v", recv)
	}

	return nil
}

func (s *GreeterServer) SayRoute(stream pb.Greeter_SayRouteServer) error {
	n := 0
	for true {
		err := stream.Send(&pb.HelloReply{Message: fmt.Sprintf("say.route.%d", n)})
		if err != nil {
			return err
		}

		recv, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		n++
		log.Printf("recv: %v", recv)
	}

	return nil
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis, _ := net.Listen("tcp", ":"+port)
	server.Serve(lis)
}
