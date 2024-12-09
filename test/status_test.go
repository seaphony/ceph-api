package test

import (
	"testing"
	"time"

	pb "github.com/clyso/ceph-api/api/gen/grpc/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_List_Ceph_Status(t *testing.T) {
	r := require.New(t)
	client := pb.NewStatusClient(admConn)
	res, err := client.GetCephStatus(tstCtx, &emptypb.Empty{})
	r.NoError(err)
	r.NotEmpty(res)
	r.NotEmpty(res.Fsid)
	r.NotEmpty(res.QuorumAge)
	r.NotEmpty(res.Health)
	r.NotEmpty(res.Osdmap)
	r.NotEmpty(res.Mgrmap)
	r.NotEmpty(res.Monmap)

}

func Test_GetCephMonDump(t *testing.T) {
	r := require.New(t)
	client := pb.NewStatusClient(admConn)

	res, err := client.GetCephMonDump(tstCtx, &emptypb.Empty{})

	r.NoError(err)
	r.NotNil(res)

	// Validate required fields in CephMonDumpResponse
	r.NotEmpty(res.Epoch, "Epoch should not be empty")
	r.NotEmpty(res.Fsid, "Fsid should not be empty")
	r.NotEmpty(res.Modified, "Modified timestamp should not be empty")
	r.NotEmpty(res.Created, "Created timestamp should not be empty")
	// created time after 1 january 2024
	r.True(res.Created.AsTime().After(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)), "Created timestamp should be after 1 January 2024")
	// created time should be less than current time
	r.True(res.Created.AsTime().Before(time.Now()), "Created timestamp should be before current time")
	// modified time should be greater than or equal to created time
	r.True((res.Modified.AsTime().After(res.Created.AsTime()) || res.Modified.AsTime().Equal(res.Created.AsTime())), "Modified timestamp should be after Created timestamp")

	r.NotZero(res.MinMonRelease, "MinMonRelease should not be zero")
	r.NotEmpty(res.MinMonReleaseName, "MinMonReleaseName should not be empty")
	r.NotNil(res.Features, "Features should not be nil")
	r.NotNil(res.Mons, "Mons should not be nil")
	r.NotEmpty(res.Quorum, "Quorum should not be empty")

	// Validate the Features sub-message
	r.NotNil(res.Features.Persistent, "Features.Persistent should not be nil")

	// Validate the Mons repeated field
	for _, mon := range res.Mons {
		r.NotNil(mon.Rank, "Mon.Rank should not be zero")
		r.NotEmpty(mon.Name, "Mon.Name should not be empty")
		r.NotEmpty(mon.Addr, "Mon.Addr should not be empty")
		r.NotEmpty(mon.PublicAddr, "Mon.PublicAddr should not be empty")
		r.NotNil(mon.PublicAddrs, "Mon.PublicAddrs should not be nil")

		// Validate PublicAddrs repeated field
		for _, addrVec := range mon.PublicAddrs.Addrvec {
			r.NotEmpty(addrVec.Type, "AddrVec.Type should not be empty")
			r.NotEmpty(addrVec.Addr, "AddrVec.Addr should not be empty")
		}
	}
}
