package main

import (
	"context"
	"log"

	pb "github.com/knight-zlm/tag-service/proto"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	clientConn, err := GetClientConn(ctx, "localhost:8004", nil)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, err := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})
	if err != nil {
		log.Fatalf("tagServiceClient.GetTagList err:%v\n", err)
	}
	log.Printf("resp:%v\n", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}
