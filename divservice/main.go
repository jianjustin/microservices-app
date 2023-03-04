package main

import (
	"context"
	"fmt"
	"go-microservices/divservice/config"
	etcd_pkg "go-microservices/divservice/pkg"
	pb "go-microservices/divservice/proto"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

type server struct {
	Addr string
	pb.UnimplementedDivServiceServer
}

func (s *server) HandleDiv(ctx context.Context, in *pb.DivRequest) (*pb.DivReply, error) {
	log.Printf("%s/Received: %f / %f", s.Addr, in.A, in.B)
	res := in.A + in.B
	return &pb.DivReply{Result: res}, nil
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
			pb.RegisterDivServiceServer(s, &server{Addr: fmt.Sprintf("0.0.0.0%s", addr)})
			log.Printf("server listening at %v", lis.Addr())
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}(addr)

	}
	wg.Wait()
}
