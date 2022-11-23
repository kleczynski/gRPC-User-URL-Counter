package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	pb "rectask/gen/proto"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterRecTaskServer(s, &RecTask{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(ctx context.Context, str string) (net.Conn, error) {
	return lis.Dial()
}

func TestAuth(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			t.Fatalf("Err: %v", err)
		}
	}(conn)
	client := pb.NewRecTaskClient(conn)
	request := empty.Empty{}
	resp, err := client.Auth(ctx, &request)
	if err != nil {
		t.Fatalf("Auth failed: %v", err)
	}
	if len(resp.Token) != 10 {
		t.Fatalf("Failed generated token: %v", resp.Token)
	}

}

//func TestAllPath(t *testing.T) {
//	ctx := context.Background()
//	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		t.Fatalf("Failed to dial bufnet: %v", err)
//	}
//	defer func(conn *grpc.ClientConn) {
//		err := conn.Close()
//		if err != nil {
//			t.Fatalf("Err: %v", err)
//		}
//	}(conn)
//	client := pb.NewRecTaskClient(conn)
//	md := metadata.New(map[string]string{"Authorization": "ZXCVBNMASD", "grpcgateway-http-path": "/whatever"})
//	if len(md.Get("Authorization")) != 1 {
//		t.Fatalf("Invalid Authorization token: %v", md.Get("Authorization")[0])
//	}
//	if len(md.Get("grpcgateway-http-path")) != 1 {
//		t.Fatalf("Invalid path: %v", md.Get("grpcgateway-http-path")[0])
//	}
//	request := empty.Empty{}
//	resp, err := client.AllPath(ctx, &request)
//	if err != nil {
//		t.Fatalf("AllPath failed: %v", err)
//	}
//
//	fmt.Printf("Response: %v", resp)
//
//}
