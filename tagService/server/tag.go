package server

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/metadata"

	"github.com/knight-zlm/tag-service/pkg/bapi"
	"github.com/knight-zlm/tag-service/pkg/errcode"
	pb "github.com/knight-zlm/tag-service/proto"
)

type Auth struct {
	AppKey    string
	AppSecret string
}

type TagServer struct {
	auth *Auth
}

func (a *Auth) GetAppKey() string {
	return "go_tour"
}

func (a *Auth) GetAppSecret() string {
	return "zlm"
}

func (a *Auth) Check(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	var appKey, appSecret string
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}

	return nil
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	if err := t.auth.Check(ctx); err != nil {
		return nil, err
	}
	api := bapi.NewAPI("http://127.0.0.1:8008")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFall)
	}

	tagList := pb.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, err
	}

	return &tagList, nil
}
