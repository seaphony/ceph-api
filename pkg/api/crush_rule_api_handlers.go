package api

import (
	"context"
	"encoding/json"

	pb "github.com/seaphony/ceph-api/api/gen/grpc/go"
	"github.com/seaphony/ceph-api/pkg/rados"
	"github.com/seaphony/ceph-api/pkg/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewCrushRuleAPI(radosSvc *rados.Svc) pb.CrushRuleServer {
	return &crushRuleAPI{
		radosSvc: radosSvc,
	}
}

type crushRuleAPI struct {
	radosSvc *rados.Svc
}

type crushDump struct {
	Rules []*pb.Rule `json:"rules"`
}

// CreateRule implements pb.CrushRuleServer.
func (c *crushRuleAPI) CreateRule(ctx context.Context, req *pb.CreateRuleRequest) (*emptypb.Empty, error) {
	if err := user.HasPermissions(ctx, user.ScopeOsd, user.PermCreate); err != nil {
		return nil, err
	}

	var cmdMap map[string]interface{}

	if req.PoolType == "erasure" {
		cmdMap = map[string]interface{}{
			"prefix":  "osd crush rule create-erasure",
			"name":    req.Name,
			"profile": req.Profile,
			"type":    req.FailureDomain,
			"class":   req.DeviceClass,
			"format":  "json",
		}
	} else {
		cmdMap = map[string]interface{}{
			"prefix": "osd crush rule create-replicated",
			"name":   req.Name,
			"root":   req.Root,
			"type":   req.FailureDomain,
			"class":  req.DeviceClass,
			"format": "json",
		}
	}

	cmdBytes, err := json.Marshal(cmdMap)
	if err != nil {
		return nil, err
	}

	_, err = c.radosSvc.ExecMon(ctx, string(cmdBytes))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// DeleteRule implements pb.CrushRuleServer.
func (c *crushRuleAPI) DeleteRule(ctx context.Context, req *pb.DeleteRuleRequest) (*emptypb.Empty, error) {
	if err := user.HasPermissions(ctx, user.ScopeOsd, user.PermDelete); err != nil {
		return nil, err
	}

	cmdMap := map[string]interface{}{
		"prefix": "osd crush rule rm",
		"name":   req.Name,
		"format": "json",
	}

	cmdBytes, err := json.Marshal(cmdMap)
	if err != nil {
		return nil, err
	}

	_, err = c.radosSvc.ExecMon(ctx, string(cmdBytes))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// GetRule implements pb.CrushRuleServer.
func (c *crushRuleAPI) GetRule(ctx context.Context, req *pb.GetRuleRequest) (*pb.GetRuleResponse, error) {
	if err := user.HasPermissions(ctx, user.ScopeOsd, user.PermRead); err != nil {
		return nil, err
	}
	const cmdTempl = `{"prefix": "osd crush dump", "format": "json"}`
	res, err := c.radosSvc.ExecMon(ctx, cmdTempl)
	if err != nil {
		return nil, err
	}
	var dump crushDump
	if err := json.Unmarshal(res, &dump); err != nil {
		return nil, err
	}

	// Find the rule by name.
	for _, rule := range dump.Rules {
		if rule.RuleName == req.Name {
			return &pb.GetRuleResponse{Rule: rule}, nil
		}
	}

	return nil, nil
}

// ListRules implements pb.CrushRuleServer.
func (c *crushRuleAPI) ListRules(ctx context.Context, req *emptypb.Empty) (*pb.ListRulesResponse, error) {
	if err := user.HasPermissions(ctx, user.ScopeOsd, user.PermRead); err != nil {
		return nil, err
	}
	const cmdTempl = `{"prefix": "osd crush dump", "format": "json"}`
	res, err := c.radosSvc.ExecMon(ctx, cmdTempl)
	if err != nil {
		return nil, err
	}
	var dump crushDump
	if err := json.Unmarshal(res, &dump); err != nil {
		return nil, err
	}

	return &pb.ListRulesResponse{Rules: dump.Rules}, nil
}
