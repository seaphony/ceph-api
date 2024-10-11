package api

import (
	"context"

	pb "github.com/seaphony/ceph-api/api/gen/grpc/go"
	"github.com/seaphony/ceph-api/pkg/user"
)

func NewHostsAPI() pb.HostServer {
	return &hostsAPI{
		// svc: svc,
	}
}

type hostsAPI struct {
	// svc *user.Service
}

// GetHost implements pb.HostServer.
func (h *hostsAPI) GetHost(ctx context.Context, r *pb.GetHostReq) (*pb.HostResp, error) {
	if err := user.HasPermissions(ctx, user.ScopeHosts, user.PermRead); err != nil {
		return nil, err
	}
	const cmdTempl = `{"prefix": "auth caps", "entity": "%s", "caps": [%s]}`
	caps := make([]string, 0, len(req.Capabilities)*2)
	for k, v := range req.Capabilities {
		caps = append(caps, strconv.Quote(k), strconv.Quote(v))
	}
	monCmd := fmt.Sprintf(cmdTempl, req.UserEntity, strings.Join(caps, ","))
	_, err := c.radosSvc.ExecMon(ctx, monCmd)
	if err != nil {
		if errors.Is(err, gorados.ErrNotFound) {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	return &pb.HostResp{Name: "clyso"}, nil
}
