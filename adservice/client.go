package main

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	pb "go-microservices/adservice/proto"
)

func main() {
	reg := etcd.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"127.0.0.1:2379"}
	})

	service := micro.NewService(
		micro.Registry(reg),
	)
	service.Init()

	client := pb.NewHealthService("adservice", service.Client())
	res, err := client.Check(context.Background(), &pb.HealthCheckRequest{Service: "client"})

	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
