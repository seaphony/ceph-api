package test

import (
	"context"
	"testing"

	pb "github.com/clyso/ceph-api/api/gen/grpc/go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_ClusterStatus(t *testing.T) {
	r := require.New(t)
	client := pb.NewClusterClient(admConn)

	res, err := client.GetStatus(tstCtx, &emptypb.Empty{})
	r.NoError(err)
	initStatus := res.Status
	newStatus := pb.ClusterStatus_INSTALLED
	if initStatus == newStatus {
		newStatus = pb.ClusterStatus_POST_INSTALLED
	}

	_, err = client.UpdateStatus(tstCtx, &pb.ClusterStatus{Status: newStatus})
	r.NoError(err)
	t.Cleanup(func() {
		client.UpdateStatus(tstCtx, &pb.ClusterStatus{Status: initStatus})
	})

	res, err = client.GetStatus(tstCtx, &emptypb.Empty{})
	r.NoError(err)
	r.EqualValues(newStatus, res.Status)
}

func Test_ClusterUsers(t *testing.T) {
	r := require.New(t)
	client := pb.NewClusterClient(admConn)
	const (
		user = "client.test"
	)

	users, err := client.GetUsers(tstCtx, &emptypb.Empty{})
	r.NoError(err, "get all users")

	_, err = client.CreateUser(tstCtx, &pb.CreateClusterUserReq{
		Capabilities: map[string]string{"mon": "allow r"},
		UserEntity:   user,
	})
	r.NoError(err, "create a new test user %s", user)
	t.Cleanup(func() {
		// delete test user on exit
		client.DeleteUser(context.Background(), &pb.DeleteClusterUserReq{UserEntity: user})
	})

	users2, err := client.GetUsers(tstCtx, &emptypb.Empty{})
	r.NoError(err, "get all users including a new one")
	r.EqualValues(len(users.Users)+1, len(users2.Users), "users number increased")
	var created *pb.ClusterUser = nil
	for i, v := range users2.Users {
		if v.Entity == user {
			created = users2.Users[i]
			break
		}
	}
	r.NotNil(created, "new user created")
	r.Len(created.Caps, 1, "new user has correct capabilities")
	r.EqualValues(created.Caps["mon"], "allow r", "new user has correct capabilities")

	exp, err := client.ExportUser(tstCtx, &pb.ExportClusterUserReq{Entities: []string{user}})
	r.NoError(err, "new user can be exported")
	r.Contains(string(exp.Data), `mon = "allow r"`, "new user export conains correct caps")

	_, err = client.UpdateUser(tstCtx, &pb.UpdateClusterUserReq{
		UserEntity:   user,
		Capabilities: map[string]string{"mon": "allow w"}})
	r.NoError(err, "new user caps updated")

	exp, err = client.ExportUser(tstCtx, &pb.ExportClusterUserReq{Entities: []string{user}})
	r.NoError(err)
	r.Contains(string(exp.Data), `mon = "allow w"`, "export contains updated caps")
	r.NotContains(string(exp.Data), `mon = "allow r"`, "export does not contains old caps")

	users2, err = client.GetUsers(tstCtx, &emptypb.Empty{})
	r.NoError(err)
	r.NotEmpty(users2.Users)
	r.EqualValues(len(users.Users)+1, len(users2.Users))
	created = nil
	for i, v := range users2.Users {
		if v.Entity == user {
			created = users2.Users[i]
			break
		}
	}
	r.NotNil(created, "list user returns updated caps")
	r.Len(created.Caps, 1, "list user returns updated caps")
	r.EqualValues(created.Caps["mon"], "allow w", "list user returns updated caps")

	_, err = client.DeleteUser(tstCtx, &pb.DeleteClusterUserReq{
		UserEntity: user})
	r.NoError(err, "delete new user")

	users2, err = client.GetUsers(tstCtx, &emptypb.Empty{})
	r.NoError(err)
	r.EqualValues(len(users.Users), len(users2.Users), "user was removed from list")
	created = nil
	for i, v := range users2.Users {
		if v.Entity == user {
			created = users2.Users[i]
			break
		}
	}
	r.Nil(created, "user was removed from list")

	_, err = client.CreateUser(tstCtx, &pb.CreateClusterUserReq{ImportData: exp.Data})
	r.NoError(err, "user was imported back from export data")

	users2, err = client.GetUsers(tstCtx, &emptypb.Empty{})
	r.NoError(err)
	r.EqualValues(len(users.Users)+1, len(users2.Users), "user is back after import")
	created = nil
	for i, v := range users2.Users {
		if v.Entity == user {
			created = users2.Users[i]
			break
		}
	}
	r.NotNil(created, "user is back after import")
}
