package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-microservices/subservice/config"
	etcd_pkg "go-microservices/subservice/pkg"
	pb "go-microservices/subservice/proto"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type server struct {
	Addr string
	pb.UnimplementedSubServiceServer
}

func (s *server) HandleSub(ctx context.Context, in *pb.SubRequest) (*pb.SubReply, error) {
	log.Printf("%s/Received: %f - %f", s.Addr, in.A, in.B)
	res := in.A - in.B
	return &pb.SubReply{Result: res}, nil
}

func main() {

	client, err := clientv3.New(clientv3.Config{Endpoints: []string{"0.0.0.0:2379"}, DialTimeout: time.Second * 5})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for _, addr := range config.Addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()

			go func(addr string) {
				err := etcd_pkg.Register(ctx, client, config.ServiceName, "localhost"+addr)
				if err != nil {
					log.Fatal(err)
				}
			}(addr)

			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer()
			pb.RegisterSubServiceServer(s, &server{Addr: fmt.Sprintf("0.0.0.0%s", addr)})
			log.Printf("server listening at %v", lis.Addr())
			go func() {
				log.Fatalln(s.Serve(lis))
			}()

			conn, err := grpc.DialContext(
				context.Background(),
				fmt.Sprintf("0.0.0.0%s", addr),
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				log.Fatalln("Failed to dial server:", err)
			}

			gwmux := runtime.NewServeMux()
			// Register Greeter
			err = pb.RegisterSubServiceHandler(context.Background(), gwmux, conn)
			if err != nil {
				log.Fatalln("Failed to register gateway:", err)
			}

			gwServer := &http.Server{
				Addr:    config.HTTP_ADDR[addr],
				Handler: gwmux,
			}

			log.Println("Serving gRPC-Gateway on http://0.0.0.0" + config.HTTP_ADDR[addr])
			log.Fatalln(gwServer.ListenAndServe())
		}(addr)

	}
	wg.Wait()
}
