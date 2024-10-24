package test

import (
	"testing"

	pb "github.com/seaphony/ceph-api/api/gen/grpc/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_ListRlues(t *testing.T) {
	r := require.New(t)
	client := pb.NewCrushRuleClient(admConn)

	res, err := client.ListRules(tstCtx, &emptypb.Empty{})
	r.NoError(err)
	rules := res.Rules
	r.NotEmpty(rules)
}
