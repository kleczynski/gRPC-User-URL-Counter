package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "rectask/gen/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	client := pb.NewRecTaskClient(conn)
	request := empty.Empty{}
	_, err = client.Auth(context.Background(), &request)
	if err != nil {
		log.Println(err)
	}
}
