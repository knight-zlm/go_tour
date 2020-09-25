package main

import (
	"context"
	"flag"
	"fmt"
	pb "grpc-demo/proto"
	"io"
	"log"

	"google.golang.org/grpc"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8008", "启动端口号")
	flag.Parse()
}

func main() {
	conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	err := SayHello(client)
	if err != nil {
		return
	}
}

func SayHello(client pb.GreeterClient) error {
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "zlm"})
	if err != nil {
		return err
	}
	fmt.Printf("client.SayHellow resp: %s\n", reply.Message)
	return nil
}

func SayList(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayList(context.Background(), r)
	if err != nil {
		return err
	}
	for true {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: %v", recv)
	}
	return nil
}

func SayRecord(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayRecord(context.Background())
	if err != nil {
		return err
	}

	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}

	recv, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("resp err: %v", recv)
	return nil
}

func SayRoute(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, err := client.SayRoute(context.Background())
	if err != nil {
		return err
	}

	for n := 0; n < 6; n++ {
		err := stream.Send(&pb.HelloRequest{Name: fmt.Sprintf("router.%d", n)})
		if err != nil {
			return err
		}

		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp err: %v", recv)
	}

	err = stream.CloseSend()
	if err != nil {
		return err
	}

	return nil
}
