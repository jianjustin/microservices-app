package service

import (
	"context"
	"github.com/google/uuid"
	"go-microservices/resourceservice/model"
	pb "go-microservices/resourceservice/proto"
	"log"
)

type ResourceServer struct {
	Addr string
	pb.UnimplementedResourceServiceServer
}

func (s *ResourceServer) AddResource(ctx context.Context, in *pb.AddResourceRequest) (*pb.AddResourceReply, error) {
	log.Printf("%s/Received", s.Addr)

	item := model.Resource{
		Id:           uuid.New().String(),
		ParentId:     in.ParentId,
		ResourceType: in.ResourceType,
		Name:         in.Name,
		Content:      in.Content,
		StdDelete:    0,
	}

	if model.DB == nil {
		model.DB = &[]model.Resource{}
	}

	*model.DB = append(*model.DB, item)

	return &pb.AddResourceReply{
		Code: 200,
		Msg:  "添加成功",
		Data: &pb.ResourceReply{
			Id:            item.Id,
			ParentId:      item.ParentId,
			ResourceType:  item.ResourceType,
			Name:          item.Name,
			Content:       item.Content,
			LastUpdatedAt: nil,
		},
	}, nil
}

func (s *ResourceServer) EditResource(ctx context.Context, in *pb.UpdateResourceRequest) (*pb.UpdateResourceReply, error) {
	for _, resource := range *model.DB {
		if resource.Id == in.Id {
			resource.Name = in.Data.Name
			return &pb.UpdateResourceReply{
				Code: 200,
				Msg:  "修改成功",
				Data: &pb.ResourceReply{
					Id:           resource.Id,
					ParentId:     resource.ParentId,
					ResourceType: resource.ResourceType,
					Name:         resource.Name,
					Content:      resource.Content,
				},
			}, nil
		}
	}

	return nil, nil
}
func (s *ResourceServer) DeleteResource(ctx context.Context, in *pb.DeleteResourceRequest) (*pb.DeleteResourceReply, error) {

	db := &[]model.Resource{}
	for _, resource := range *model.DB {
		if resource.Id != in.Id {
			*db = append(*db, resource)
		}
	}

	model.DB = db
	return &pb.DeleteResourceReply{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}
func (s *ResourceServer) GetResourceById(ctx context.Context, in *pb.GetOneResourceRequest) (*pb.GetOneResourceReply, error) {
	for _, resource := range *model.DB {
		if resource.Id == in.Id {
			return &pb.GetOneResourceReply{
				Code: 0,
				Data: &pb.ResourceReply{
					Id:           resource.Id,
					ParentId:     resource.ParentId,
					ResourceType: resource.ResourceType,
					Name:         resource.Name,
					Content:      resource.Content,
				},
			}, nil
		}
	}

	return nil, nil
}
func (s *ResourceServer) GetResourceByPage(ctx context.Context, in *pb.GetPageResourceRequest) (*pb.GetPageResourceReply, error) {
	start := (in.Page - 1) * in.PageSize
	end := start + in.PageSize

	res := []*pb.ResourceReply{}
	data := (*model.DB)[start:end]
	for _, resource := range data {
		res = append(res, &pb.ResourceReply{
			Id:           resource.Id,
			ParentId:     resource.ParentId,
			ResourceType: resource.ResourceType,
			Name:         resource.Name,
			Content:      resource.Content,
		})
	}

	return &pb.GetPageResourceReply{
		List: res,
		Page: &pb.PageInfoReply{
			Total:    int64(len(*model.DB)),
			Page:     in.Page,
			PageSize: in.PageSize,
		},
		Code: 200,
		Msg:  "查询成功",
	}, nil
}
func (s *ResourceServer) GetAllResources(ctx context.Context, in *pb.GetAllResourceRequest) (*pb.GetAllResourceReply, error) {
	res := []*pb.ResourceReply{}
	for _, resource := range *model.DB {
		res = append(res, &pb.ResourceReply{
			Id:           resource.Id,
			ParentId:     resource.ParentId,
			ResourceType: resource.ResourceType,
			Name:         resource.Name,
			Content:      resource.Content,
		})
	}

	return &pb.GetAllResourceReply{
		List: res,
		Code: 200,
		Msg:  "查询成功",
	}, nil
}
