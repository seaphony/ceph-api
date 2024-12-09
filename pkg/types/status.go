package types

import (
	"time"

	pb "github.com/clyso/ceph-api/api/gen/grpc/go"
)

type CephMonDumpResponse struct {
	Epoch             int32                    `json:"epoch,omitempty"`
	Fsid              string                   `json:"fsid"` // required
	Modified          *time.Time               `json:"modified,omitempty"`
	Created           *time.Time               `json:"created"` // required
	MinMonRelease     int32                    `json:"min_mon_release,omitempty"`
	MinMonReleaseName string                   `json:"min_mon_release_name,omitempty"`
	ElectionStrategy  int32                    `json:"election_strategy,omitempty"`
	DisallowedLeaders string                   `json:"disallowed_leaders,omitempty"`
	StretchMode       bool                     `json:"stretch_mode,omitempty"`
	TiebreakerMon     string                   `json:"tiebreaker_mon,omitempty"`
	RemovedRanks      string                   `json:"removed_ranks,omitempty"`
	Features          *pb.CephMonDumpFeatures  `json:"features,omitempty"`
	Mons              []*pb.CephMonDumpMonInfo `json:"mons,omitempty"`
	Quorum            []int32                  `json:"quorum,omitempty"`
}
