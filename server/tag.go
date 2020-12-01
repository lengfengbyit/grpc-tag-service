package server

import (
	"context"
	"encoding/json"
	"go-tour/grpc-tag-service/pkg/bapi"
	pb "go-tour/grpc-tag-service/proto"
)

const API_DOMAIN = "http://127.0.0.1:8000"

type TagServer struct{}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	api := bapi.NewAPI(API_DOMAIN)
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, err
	}

	tagList := &pb.GetTagListReply{}
	err = json.Unmarshal(body, tagList)
	if err != nil {
		return nil, err
	}
	return tagList, nil
}
