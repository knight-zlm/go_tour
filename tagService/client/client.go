package main

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"

	"github.com/knight-zlm/tag-service/internal/middleware"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc/codes"

	pb "github.com/knight-zlm/tag-service/proto"
	"google.golang.org/grpc"
)

type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func main() {
	auth := Auth{
		AppKey:    "go_tour",
		AppSecret: "zlm",
	}
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithPerRPCCredentials(&auth)}
	md := metadata.New(map[string]string{"go": "programming", "tour": "book"})
	newCtx := metadata.NewOutgoingContext(ctx, md)
	clientConn, err := GetClientConn(newCtx, "localhost:8004", opts)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err:%v\n", err)
	}
	log.Printf("resp:%v\n", resp)
}

func main1() {
	ctx := context.Background()
	newCtx := metadata.AppendToOutgoingContext(ctx, "zlm", "go编程之旅")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpcMiddleware.ChainUnaryClient(
			grpcRetry.UnaryClientInterceptor(
				grpcRetry.WithMax(2),
				grpcRetry.WithCodes(
					codes.Unknown,
					codes.Internal,
					codes.DeadlineExceeded,
				),
			),
			middleware.UnaryContextTimeout(),
		),
	))
	clientConn, err := GetClientConn(newCtx, "localhost:8004", opts)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err:%v\n", err)
	}
	log.Printf("resp:%v\n", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}
