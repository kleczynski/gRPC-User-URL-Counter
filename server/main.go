package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"net"
	"net/http"
	pb "rectask/gen/proto"
	"strings"
)

type RecTask struct {
	pb.UnimplementedRecTaskServer
}

// global variable that store map in structure like: { "user_toke": counter }
var tokens []string
var paths map[[2]string]int32

func (s *RecTask) Auth(ctx context.Context, empty *empty.Empty) (*pb.TokenRender, error) {
	// When user hit /auth endpoint, then token is render
	response := &pb.TokenRender{
		Token: randStringBytes(10),
	}
	// Appending token into global variable
	tokens = append(tokens, response.Token)
	return response, nil
}

func (s *RecTask) AllPath(ctx context.Context, empty *empty.Empty) (*empty.Empty, error) {
	// Checking if paths map is empty or not
	if paths == nil {
		paths = make(map[[2]string]int32)
	}

	// Extracting data from meta and getting Authorization token
	md, _ := metadata.FromIncomingContext(ctx)
	var token string
	var newToken []string
	if len(md.Get("Authorization")) != 1 {
		err := status.Error(http.StatusUnauthorized, "U cannot enter this endpoint, first generate token in /auth")
		if err != nil {
			return nil, err
		}
	} else {
		token = md.Get("Authorization")[0]
		newToken = strings.SplitAfter(token, " ")
	}
	// Extracting current path from metadata
	var path string
	//Checking if user made right request, bad request = "/", good one = "/."
	if len(md.Get("grpcgateway-http-path")[0]) == 1 {
		err := status.Error(http.StatusBadRequest, "Bad request, used any path in request /...")
		if err != nil {
			return nil, err
		}
	} else {
		path = md["grpcgateway-http-path"][0]
	}
	log.Println(len(path))
	pathAndToken := [2]string{path, token}
	// Checking if authorization header is in generated token
	if contains(tokens, newToken[1]) {
		// if is, then checking if array ['path', 'token'] exist, then counter increase, else counter = 1
		if _, ok := paths[pathAndToken]; ok {
			paths[pathAndToken]++
		} else {
			paths[pathAndToken] = 1
		}
		fmt.Printf("{ count: %v }\n", paths[pathAndToken])
	} else {
		// if not, then raise error with code status 401
		err := status.Error(http.StatusUnauthorized, "Invalid token: Use /auth to generate token\n")
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {
	go func() {
		// mux
		mux := runtime.NewServeMux(runtime.WithMetadata(func(_ context.Context, req *http.Request) metadata.MD {
			return metadata.New(map[string]string{
				"grpcgateway-http-path": req.URL.Path,
			})
		}))
		// register
		err := pb.RegisterRecTaskHandlerServer(context.Background(), mux, &RecTask{})
		if err != nil {
			return
		}
		// http server
		log.Println("Start listen on port :8081 - via grpc-gateway")
		log.Fatalln(http.ListenAndServe("localhost:8081", mux))
	}()
	// tcp listener, just for start server
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRecTaskServer(grpcServer, &RecTask{})
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}

}

// Helper function, first one generate random token, second one checking if string is in slice
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func contains(tokens []string, userToken string) bool {
	for _, l := range tokens {
		if l == userToken {
			return true
		}

	}
	return false
}
