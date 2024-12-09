package api

import (
	"context"
	"encoding/json"

	pb "github.com/clyso/ceph-api/api/gen/grpc/go"
	"github.com/clyso/ceph-api/pkg/rados"
	"github.com/clyso/ceph-api/pkg/types"
	"github.com/clyso/ceph-api/pkg/user"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewStatusAPI(radosSvc *rados.Svc) pb.StatusServer {
	return &statusAPI{
		radosSvc: radosSvc,
	}
}

type statusAPI struct {
	radosSvc *rados.Svc
}

func (s *statusAPI) GetCephStatus(ctx context.Context, body *emptypb.Empty) (*pb.GetCephStatusResponse, error) {
	if err := user.HasPermissions(ctx, user.ScopeMonitor, user.PermRead); err != nil {
		return nil, err
	}

	const cmdTempl = `{"prefix": "status", "format": "json"}`
	res, err := s.radosSvc.ExecMon(ctx, cmdTempl)
	if err != nil {
		return nil, err
	}
	var statusDump pb.GetCephStatusResponse
	if err := json.Unmarshal(res, &statusDump); err != nil {
		return nil, err
	}

	return &statusDump, nil
}

func (s *statusAPI) GetCephMonDump(ctx context.Context, req *emptypb.Empty) (*pb.CephMonDumpResponse, error) {
	if err := user.HasPermissions(ctx, user.ScopeMonitor, user.PermRead); err != nil {
		return nil, err
	}

	const cmdTempl = `{"prefix": "mon dump", "format": "json"}`
	res, err := s.radosSvc.ExecMon(ctx, cmdTempl)
	if err != nil {
		return nil, err
	}
	var monDump types.CephMonDumpResponse
	if err := json.Unmarshal(res, &monDump); err != nil {
		return nil, err
	}

	modifiedTimestamp := timestamppb.New(*monDump.Modified)
	createdTimestamp := timestamppb.New(*monDump.Created)
	response := pb.CephMonDumpResponse{
		Epoch:             monDump.Epoch,
		Fsid:              monDump.Fsid,
		Modified:          modifiedTimestamp,
		Created:           createdTimestamp,
		MinMonRelease:     monDump.MinMonRelease,
		MinMonReleaseName: monDump.MinMonReleaseName,
		ElectionStrategy:  monDump.ElectionStrategy,
		DisallowedLeaders: monDump.DisallowedLeaders,
		StretchMode:       monDump.StretchMode,
		TiebreakerMon:     monDump.TiebreakerMon,
		RemovedRanks:      monDump.RemovedRanks,
		Features:          monDump.Features,
		Mons:              monDump.Mons,
		Quorum:            monDump.Quorum,
	}

	return &response, nil
}
