package test

import (
	"testing"

	pb "github.com/seaphony/ceph-api/api/gen/grpc/go"
	"github.com/stretchr/testify/require"
)

func Test_HOST(t *testing.T) {
	var (
		username = "test-user"
	)
	r := require.New(t)
	client := pb.NewHostClient(admConn)

	res, err := client.GetHost(tstCtx, &pb.GetHostReq{Name: username})
	r.NoError(err)
	r.EqualValues("clyso", res.Name)
}

func Test_HOST_htttp(t *testing.T) {
	var (
		username = "test-user"
	)
	r := require.New(t)
	client := pb.NewHostClient(admConn)

	res, err := client.GetHost(tstCtx, &pb.GetHostReq{Name: username})
	r.NoError(err)
	r.EqualValues("clyso", res.Name)
}
