package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"go-microservices/resourceservice/config"
	etcd_pkg "go-microservices/resourceservice/pkg"
	pb "go-microservices/resourceservice/proto"
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
	c := pb.NewResourceServiceClient(conn)

	res, err := c.AddResource(context.Background(), &pb.AddResourceRequest{
		ParentId:     uuid.New().String(),
		ResourceType: 0,
		Name:         "resource 1",
		Content:      "",
	})
	if err != nil {
		log.Printf(err.Error())
	}

	data, err := c.GetResourceById(context.Background(), &pb.GetOneResourceRequest{Id: res.Data.Id})
	if data != nil {
		log.Printf("Name:%s,ParentId:%s", res.Data.Name, res.Data.ParentId)
	}

	res, err = c.AddResource(context.Background(), &pb.AddResourceRequest{
		ParentId:     uuid.New().String(),
		ResourceType: 0,
		Name:         "resource 2",
		Content:      "",
	})

	list, err := c.GetResourceByPage(context.Background(), &pb.GetPageResourceRequest{
		Page:     1,
		PageSize: 1,
	})
	for i, item := range list.List {
		log.Printf("%d:%s", i, item.Name)
	}

	datas, err := c.GetAllResources(context.Background(), &pb.GetAllResourceRequest{})
	for i, item := range datas.List {
		log.Printf("%d:%s", i, item.Name)
	}
}
