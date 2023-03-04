package main

import (
	"context"
	"flag"
	"fmt"
	"go-microservices/mulservice/config"
	etcd_pkg "go-microservices/mulservice/pkg"
	pb "go-microservices/mulservice/proto"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

func main() {
	flag.Parse()
	client, err := clientv3.New(clientv3.Config{Endpoints: []string{"0.0.0.0:2379"}, DialTimeout: time.Second * 5})
	if err != nil {
		log.Fatal(err)
	}

	b := etcd_pkg.NewBuilder(client)

	resolver.Register(b)

	conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", etcd_pkg.Scheme, config.ServiceName), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMulServiceClient(conn)

	i := 1

	for {
		r, err := c.HandleMul(context.Background(), &pb.MulRequest{
			A: 1,
			B: 2,
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(r.Result)
		}
		i++
		time.Sleep(time.Second)
	}
}
