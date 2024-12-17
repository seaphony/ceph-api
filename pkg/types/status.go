package types

import (
	"strings"
	"time"

	pb "github.com/clyso/ceph-api/api/gen/grpc/go"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type CephOsdDumpResponse struct {
	Epoch                  int32                                    `json:"epoch,omitempty"`
	Fsid                   string                                   `json:"fsid"`
	Created                CephTimestamp                            `json:"created"`
	Modified               CephTimestamp                            `json:"modified,omitempty"`
	LastUpChange           CephTimestamp                            `json:"last_up_change,omitempty"`
	LastInChange           CephTimestamp                            `json:"last_in_change,omitempty"`
	Flags                  string                                   `json:"flags,omitempty"`
	FlagsNum               int32                                    `json:"flags_num,omitempty"`
	FlagsSet               []string                                 `json:"flags_set,omitempty"`
	CrushVersion           int32                                    `json:"crush_version,omitempty"`
	FullRatio              float64                                  `json:"full_ratio,omitempty"`
	BackfillfullRatio      float64                                  `json:"backfillfull_ratio,omitempty"`
	NearfullRatio          float64                                  `json:"nearfull_ratio,omitempty"`
	ClusterSnapshot        string                                   `json:"cluster_snapshot,omitempty"`
	PoolMax                int32                                    `json:"pool_max,omitempty"`
	MaxOsd                 int32                                    `json:"max_osd,omitempty"`
	RequireMinCompatClient string                                   `json:"require_min_compat_client,omitempty"`
	MinCompatClient        string                                   `json:"min_compat_client,omitempty"`
	RequireOsdRelease      string                                   `json:"require_osd_release,omitempty"`
	AllowCrimson           bool                                     `json:"allow_crimson,omitempty"`
	Pools                  []OsdDumpPool                            `json:"pools,omitempty"`
	Osds                   []*pb.OsdDumpOsdInfo                     `json:"osds,omitempty"`
	OsdXinfo               []*OsdDumpOsdXInfo                       `json:"osd_xinfo,omitempty"`
	PgUpmap                []*structpb.Value                        `json:"pg_upmap,omitempty"`
	PgUpmapItems           []*structpb.Value                        `json:"pg_upmap_items,omitempty"`
	PgUpmapPrimaries       []*structpb.Value                        `json:"pg_upmap_primaries,omitempty"`
	PgTemp                 []*structpb.Value                        `json:"pg_temp,omitempty"`
	PrimaryTemp            []*structpb.Value                        `json:"primary_temp,omitempty"`
	Blocklist              map[string]CephTimestamp                 `json:"blocklist,omitempty"`
	RangeBlocklist         *structpb.Struct                         `json:"range_blocklist,omitempty"`
	ErasureCodeProfiles    map[string]*pb.OsdDumpErasureCodeProfile `json:"erasure_code_profiles,omitempty"`
	RemovedSnapsQueue      []*structpb.Value                        `json:"removed_snaps_queue,omitempty"`
	NewRemovedSnaps        []*structpb.Value                        `json:"new_removed_snaps,omitempty"`
	NewPurgedSnaps         []*structpb.Value                        `json:"new_purged_snaps,omitempty"`
	CrushNodeFlags         *structpb.Struct                         `json:"crush_node_flags,omitempty"`
	DeviceClassFlags       *structpb.Struct                         `json:"device_class_flags,omitempty"`
	StretchMode            *pb.OsdDumpStretchMode                   `json:"stretch_mode,omitempty"`
}

type OsdDumpPool struct {
	Pool                              int32                      `json:"pool,omitempty"`
	PoolName                          string                     `json:"pool_name,omitempty"`
	CreateTime                        CephTimestamp              `json:"create_time"`
	Flags                             int64                      `json:"flags,omitempty"`
	FlagsNames                        string                     `json:"flags_names,omitempty"`
	Type                              int32                      `json:"type,omitempty"`
	Size                              int32                      `json:"size,omitempty"`
	MinSize                           int32                      `json:"min_size,omitempty"`
	CrushRule                         int32                      `json:"crush_rule,omitempty"`
	PeeringCrushBucketCount           int32                      `json:"peering_crush_bucket_count,omitempty"`
	PeeringCrushBucketTarget          int32                      `json:"peering_crush_bucket_target,omitempty"`
	PeeringCrushBucketBarrier         int32                      `json:"peering_crush_bucket_barrier,omitempty"`
	PeeringCrushBucketMandatoryMember int32                      `json:"peering_crush_bucket_mandatory_member,omitempty"`
	ObjectHash                        int32                      `json:"object_hash,omitempty"`
	PgAutoscaleMode                   string                     `json:"pg_autoscale_mode,omitempty"`
	PgNum                             int32                      `json:"pg_num,omitempty"`
	PgPlacementNum                    int32                      `json:"pg_placement_num,omitempty"`
	PgPlacementNumTarget              int32                      `json:"pg_placement_num_target,omitempty"`
	PgNumTarget                       int32                      `json:"pg_num_target,omitempty"`
	PgNumPending                      int32                      `json:"pg_num_pending,omitempty"`
	LastPgMergeMeta                   *pb.OsdDumpLastPgMergeMeta `json:"last_pg_merge_meta,omitempty"`
	LastChange                        string                     `json:"last_change,omitempty"`
	LastForceOpResend                 string                     `json:"last_force_op_resend,omitempty"`
	LastForceOpResendPrenautilus      string                     `json:"last_force_op_resend_prenautilus,omitempty"`
	LastForceOpResendPreluminous      string                     `json:"last_force_op_resend_preluminous,omitempty"`
	Auid                              uint64                     `json:"auid,omitempty"`
	SnapMode                          string                     `json:"snap_mode,omitempty"`
	SnapSeq                           uint64                     `json:"snap_seq,omitempty"`
	SnapEpoch                         uint64                     `json:"snap_epoch,omitempty"`
	PoolSnaps                         []*structpb.Value          `json:"pool_snaps,omitempty"`
	RemovedSnaps                      string                     `json:"removed_snaps,omitempty"`
	QuotaMaxBytes                     uint64                     `json:"quota_max_bytes,omitempty"`
	QuotaMaxObjects                   uint64                     `json:"quota_max_objects,omitempty"`
	Tiers                             []int32                    `json:"tiers,omitempty"`
	TierOf                            int32                      `json:"tier_of,omitempty"`
	ReadTier                          int32                      `json:"read_tier,omitempty"`
	WriteTier                         int32                      `json:"write_tier,omitempty"`
	CacheMode                         string                     `json:"cache_mode,omitempty"`
	TargetMaxBytes                    uint64                     `json:"target_max_bytes,omitempty"`
	TargetMaxObjects                  uint64                     `json:"target_max_objects,omitempty"`
	CacheTargetDirtyRatioMicro        uint64                     `json:"cache_target_dirty_ratio_micro,omitempty"`
	CacheTargetDirtyHighRatioMicro    uint64                     `json:"cache_target_dirty_high_ratio_micro,omitempty"`
	CacheTargetFullRatioMicro         uint64                     `json:"cache_target_full_ratio_micro,omitempty"`
	CacheMinFlushAge                  uint64                     `json:"cache_min_flush_age,omitempty"`
	CacheMinEvictAge                  uint64                     `json:"cache_min_evict_age,omitempty"`
	ErasureCodeProfile                string                     `json:"erasure_code_profile,omitempty"`
	HitSetParams                      *pb.OsdDumpHitSetParams    `json:"hit_set_params,omitempty"`
	HitSetPeriod                      uint64                     `json:"hit_set_period,omitempty"`
	HitSetCount                       uint64                     `json:"hit_set_count,omitempty"`
	UseGmtHitset                      bool                       `json:"use_gmt_hitset,omitempty"`
	MinReadRecencyForPromote          uint64                     `json:"min_read_recency_for_promote,omitempty"`
	MinWriteRecencyForPromote         uint64                     `json:"min_write_recency_for_promote,omitempty"`
	HitSetGradeDecayRate              uint64                     `json:"hit_set_grade_decay_rate,omitempty"`
	HitSetSearchLastN                 uint64                     `json:"hit_set_search_last_n,omitempty"`
	GradeTable                        []*structpb.Value          `json:"grade_table,omitempty"`
	StripeWidth                       uint64                     `json:"stripe_width,omitempty"`
	ExpectedNumObjects                uint64                     `json:"expected_num_objects,omitempty"`
	FastRead                          bool                       `json:"fast_read,omitempty"`
	Options                           *structpb.Struct           `json:"options,omitempty"`
	ApplicationMetadata               *structpb.Struct           `json:"application_metadata,omitempty"`
	ReadBalance                       *pb.OsdDumpReadBalance     `json:"read_balance,omitempty"`
}

type OsdDumpOsdXInfo struct {
	Osd                  int32         `json:"osd,omitempty"`
	DownStamp            CephTimestamp `json:"down_stamp,omitempty"`
	LaggyProbability     float64       `json:"laggy_probability,omitempty"`
	LaggyInterval        float64       `json:"laggy_interval,omitempty"`
	Features             uint64        `json:"features,omitempty"`
	OldWeight            float64       `json:"old_weight,omitempty"`
	LastPurgedSnapsScrub CephTimestamp `json:"last_purged_snaps_scrub,omitempty"`
	DeadEpoch            int32         `json:"dead_epoch,omitempty"`
}

type CephTimestamp struct {
	*timestamppb.Timestamp
}

const customTimeLayout = "2006-01-02T15:04:05.000000-0700"

// custom unmashal function for CephTimestamp
func (ct *CephTimestamp) UnmarshalJSON(data []byte) error {
	// data is a JSON string (e.g., "\"2023-05-01T12:34:56.000000-0700\"")

	// First, trim surrounding quotes.
	s := strings.Trim(string(data), `"`)

	// Handle the "0.000000" or empty-string case
	if s == "0.000000" || s == "" {
		ct.Timestamp = timestamppb.New(time.Time{})
		return nil
	}

	parsed, err := time.Parse(customTimeLayout, s)
	if err != nil {
		return err
	}
	ct.Timestamp = timestamppb.New(parsed)
	return nil
}
