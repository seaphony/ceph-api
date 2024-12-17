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

func (s *statusAPI) GetCephOsdDump(ctx context.Context, body *emptypb.Empty) (*pb.GetCephOsdDumpResponse, error) {
	if err := user.HasPermissions(ctx, user.ScopeOsd, user.PermRead); err != nil {
		return nil, err
	}

	const cmdTempl = `{"prefix": "osd dump", "format": "json"}`
	res, err := s.radosSvc.ExecMon(ctx, cmdTempl)
	if err != nil {
		return nil, err
	}

	var osdDump types.CephOsdDumpResponse
	if err := json.Unmarshal(res, &osdDump); err != nil {
		return nil, err
	}

	response := convertToPbGetCephOsdDumpResponse(osdDump)

	return response, nil
}

func convertToPbGetCephOsdDumpResponse(osdDump types.CephOsdDumpResponse) *pb.GetCephOsdDumpResponse {
	// Convert pools
	var osdDumpPools []*pb.OsdDumpPool
	for _, pool := range osdDump.Pools {
		osdDumpPools = append(osdDumpPools, &pb.OsdDumpPool{
			Pool:                              pool.Pool,
			PoolName:                          pool.PoolName,
			CreateTime:                        pool.CreateTime.Timestamp,
			Flags:                             pool.Flags,
			FlagsNames:                        pool.FlagsNames,
			Type:                              pool.Type,
			Size:                              pool.Size,
			MinSize:                           pool.MinSize,
			CrushRule:                         pool.CrushRule,
			PeeringCrushBucketCount:           pool.PeeringCrushBucketCount,
			PeeringCrushBucketTarget:          pool.PeeringCrushBucketTarget,
			PeeringCrushBucketBarrier:         pool.PeeringCrushBucketBarrier,
			PeeringCrushBucketMandatoryMember: pool.PeeringCrushBucketMandatoryMember,
			ObjectHash:                        pool.ObjectHash,
			PgAutoscaleMode:                   pool.PgAutoscaleMode,
			PgNum:                             pool.PgNum,
			PgPlacementNum:                    pool.PgPlacementNum,
			PgPlacementNumTarget:              pool.PgPlacementNumTarget,
			PgNumTarget:                       pool.PgNumTarget,
			PgNumPending:                      pool.PgNumPending,
			LastPgMergeMeta:                   pool.LastPgMergeMeta,
			LastChange:                        pool.LastChange,
			LastForceOpResend:                 pool.LastForceOpResend,
			LastForceOpResendPrenautilus:      pool.LastForceOpResendPrenautilus,
			LastForceOpResendPreluminous:      pool.LastForceOpResendPreluminous,
			Auid:                              pool.Auid,
			SnapMode:                          pool.SnapMode,
			SnapSeq:                           pool.SnapSeq,
			SnapEpoch:                         pool.SnapEpoch,
			PoolSnaps:                         pool.PoolSnaps,
			RemovedSnaps:                      pool.RemovedSnaps,
			QuotaMaxBytes:                     pool.QuotaMaxBytes,
			QuotaMaxObjects:                   pool.QuotaMaxObjects,
			Tiers:                             pool.Tiers,
			TierOf:                            pool.TierOf,
			ReadTier:                          pool.ReadTier,
			WriteTier:                         pool.WriteTier,
			CacheMode:                         pool.CacheMode,
			TargetMaxBytes:                    pool.TargetMaxBytes,
			TargetMaxObjects:                  pool.TargetMaxObjects,
			CacheTargetDirtyRatioMicro:        pool.CacheTargetDirtyRatioMicro,
			CacheTargetDirtyHighRatioMicro:    pool.CacheTargetDirtyHighRatioMicro,
			CacheTargetFullRatioMicro:         pool.CacheTargetFullRatioMicro,
			CacheMinFlushAge:                  pool.CacheMinFlushAge,
			CacheMinEvictAge:                  pool.CacheMinEvictAge,
			ErasureCodeProfile:                pool.ErasureCodeProfile,
			HitSetParams:                      pool.HitSetParams,
			HitSetPeriod:                      pool.HitSetPeriod,
			HitSetCount:                       pool.HitSetCount,
			UseGmtHitset:                      pool.UseGmtHitset,
			MinReadRecencyForPromote:          pool.MinReadRecencyForPromote,
			MinWriteRecencyForPromote:         pool.MinWriteRecencyForPromote,
			HitSetGradeDecayRate:              pool.HitSetGradeDecayRate,
			HitSetSearchLastN:                 pool.HitSetSearchLastN,
			GradeTable:                        pool.GradeTable,
			StripeWidth:                       pool.StripeWidth,
			ExpectedNumObjects:                pool.ExpectedNumObjects,
			FastRead:                          pool.FastRead,
			Options:                           pool.Options,
			ApplicationMetadata:               pool.ApplicationMetadata,
			ReadBalance:                       pool.ReadBalance,
		})
	}

	blocklistPb := make(map[string]*timestamppb.Timestamp, len(osdDump.Blocklist))
	for ip, t := range osdDump.Blocklist {
		blocklistPb[ip] = t.Timestamp
	}

	var osdXInfo []*pb.OsdDumpOsdXInfo
	for _, osdX := range osdDump.OsdXinfo {
		osdXInfo = append(osdXInfo, &pb.OsdDumpOsdXInfo{
			Osd:                  osdX.Osd,
			DownStamp:            osdX.DownStamp.Timestamp,
			LaggyProbability:     osdX.LaggyProbability,
			LaggyInterval:        osdX.LaggyInterval,
			Features:             osdX.Features,
			OldWeight:            osdX.OldWeight,
			LastPurgedSnapsScrub: osdX.LastPurgedSnapsScrub.Timestamp,
			DeadEpoch:            osdX.DeadEpoch,
		})
	}

	return &pb.GetCephOsdDumpResponse{
		Epoch:                  osdDump.Epoch,
		Fsid:                   osdDump.Fsid,
		Modified:               osdDump.Modified.Timestamp,
		Created:                osdDump.Created.Timestamp,
		LastUpChange:           osdDump.LastUpChange.Timestamp,
		LastInChange:           osdDump.LastInChange.Timestamp,
		Flags:                  osdDump.Flags,
		FlagsNum:               osdDump.FlagsNum,
		FlagsSet:               osdDump.FlagsSet,
		CrushVersion:           osdDump.CrushVersion,
		FullRatio:              osdDump.FullRatio,
		BackfillfullRatio:      osdDump.BackfillfullRatio,
		NearfullRatio:          osdDump.NearfullRatio,
		ClusterSnapshot:        osdDump.ClusterSnapshot,
		PoolMax:                osdDump.PoolMax,
		MaxOsd:                 osdDump.MaxOsd,
		RequireMinCompatClient: osdDump.RequireMinCompatClient,
		MinCompatClient:        osdDump.MinCompatClient,
		RequireOsdRelease:      osdDump.RequireOsdRelease,
		AllowCrimson:           osdDump.AllowCrimson,
		Pools:                  osdDumpPools,

		Osds:     osdDump.Osds,
		OsdXinfo: osdXInfo,

		PgUpmap:          osdDump.PgUpmap,
		PgUpmapItems:     osdDump.PgUpmapItems,
		PgUpmapPrimaries: osdDump.PgUpmapPrimaries,
		PgTemp:           osdDump.PgTemp,
		PrimaryTemp:      osdDump.PrimaryTemp,

		Blocklist:           blocklistPb,
		RangeBlocklist:      osdDump.RangeBlocklist,
		ErasureCodeProfiles: osdDump.ErasureCodeProfiles,

		RemovedSnapsQueue: osdDump.RemovedSnapsQueue,
		NewRemovedSnaps:   osdDump.NewRemovedSnaps,
		NewPurgedSnaps:    osdDump.NewPurgedSnaps,

		CrushNodeFlags:   osdDump.CrushNodeFlags,
		DeviceClassFlags: osdDump.DeviceClassFlags,
		StretchMode:      osdDump.StretchMode,
	}
}
