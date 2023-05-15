package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	calculateService "github.com/jianjustin/calculateservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

var (
	port    = flag.Int("port", 50051, "The server port")
	restful = flag.Int("restful", 8080, "the port to restful serve on")
)

type server struct {
	calculateService.UnimplementedCalculateServiceServer
}

func (s *server) Add(ctx context.Context, in *calculateService.AddRequest) (*calculateService.AddReply, error) {
	return &calculateService.AddReply{
		Res: in.A + in.B,
	}, nil
}

func (s *server) Sub(ctx context.Context, in *calculateService.SubRequest) (*calculateService.SubReply, error) {
	return &calculateService.SubReply{
		Res: in.A - in.B,
	}, nil
}

func (s *server) Mul(ctx context.Context, in *calculateService.MulRequest) (*calculateService.MulReply, error) {
	return &calculateService.MulReply{
		Res: in.A * in.B,
	}, nil
}

func (s *server) Div(ctx context.Context, in *calculateService.DivRequest) (*calculateService.DivReply, error) {
	return &calculateService.DivReply{
		Res: in.A / in.B,
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	calculateService.RegisterCalculateServiceServer(s, &server{})

	// Serve gRPC server
	log.Printf("Serving gRPC on 0.0.0.0:%d\n", *port)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("0.0.0.0:%d", *port),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = calculateService.RegisterCalculateServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *restful),
		Handler: gwmux,
	}

	log.Println(fmt.Sprintf("Serving gRPC-Gateway on http://0.0.0.0::%d", *restful))
	log.Fatalln(gwServer.ListenAndServe())

}
