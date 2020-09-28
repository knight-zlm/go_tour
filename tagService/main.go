package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/knight-zlm/tag-service/proto"
	"github.com/knight-zlm/tag-service/server"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8008", "启动端口号")
	flag.Parse()
}

//protoc --go_out=plugins=grpc:. ./proto/*.proto 编译proto
func main() {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("s.Serve err: %v", err)
	}
}
