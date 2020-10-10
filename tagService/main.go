package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/soheilhy/cmux"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/knight-zlm/tag-service/proto"
	"github.com/knight-zlm/tag-service/server"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", "8003", "启动端口号")
	flag.Parse()
}

//protoc --go_out=plugins=grpc:. ./proto/*.proto 编译proto
func main() {
	l, err := RunTcpServer(port)
	if err != nil {
		log.Fatalf("Run tcp Server err:%v", err)
	}
	m := cmux.New(l)
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())
	grpcS := RunGrpcServer(port)
	httpS := RunHttpServer(port)
	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)

	err = m.Serve()
	if err != nil {
		log.Fatalf("Run Server err:%v", err)
	}
}

func RunTcpServer(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

func RunHttpServer(port string) *http.Server {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})
	return &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
}

func RunGrpcServer(port string) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	return s
}
