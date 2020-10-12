package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

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
	flag.StringVar(&port, "port", "8004", "启动端口号")
	flag.Parse()
}

//protoc --go_out=plugins=grpc:. ./proto/*.proto 编译proto
// 添加支持JSON api的功能，需要google/api https://github.com/aspnet/AspLabs/tree/12d388c1964c8844dcbbdcd643f8bd7c6423a4c4/src/GrpcHttpApi/sample/Proto/google/api
//protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gatewat_out=logtostderr=true:. ./proto/*.proto
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

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
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
