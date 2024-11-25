package api

import (
	"context"

	pb "github.com/clyso/ceph-api/api/gen/grpc/go"
	"github.com/clyso/ceph-api/pkg/auth"
	"github.com/clyso/ceph-api/pkg/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func NewAuthAPI(svc *auth.Server) pb.AuthServer {
	return &authAPI{
		svc: svc,
	}
}

type authAPI struct {
	svc *auth.Server
}

func (a *authAPI) Check(ctx context.Context, req *pb.TokenCheckReq) (*pb.TokenCheckResp, error) {
	return nil, types.ErrNotImplemented
}

func (a *authAPI) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	res, err := a.svc.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	permissions := make(map[string]*structpb.ListValue, len(res.Permissions))
	for p, vals := range res.Permissions {
		permissions[p] = &structpb.ListValue{}
		for _, v := range vals {
			permissions[p].Values = append(permissions[p].Values, structpb.NewStringValue(v))
		}
	}
	return &pb.LoginResp{
		Token:             res.Token,
		Username:          res.User.Username,
		PwdUpdateRequired: res.User.PwdUpdateRequired,
		PwdExpirationDate: tsToPb(res.User.PwdExpirationDate),
		Sso:               false,
		Permissions:       permissions,
	}, nil
}

func (a *authAPI) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := a.svc.Logout(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
