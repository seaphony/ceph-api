package api

import (
	"context"
	"encoding/json"
	"time"

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
	if err := user.HasPermissions(ctx, user.ScopeMonitor, user.PermRead); err != nil {
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

	modifiedGoParsedTimestamp, err := parseCustomTimestamp(osdDump.Modified)
	if err != nil {
		return nil, err
	}
	modifiedTimestamp := timestamppb.New(modifiedGoParsedTimestamp)

	createdGoParsedTimestamp, err := parseCustomTimestamp(osdDump.Created)
	if err != nil {
		return nil, err
	}
	createdTimestamp := timestamppb.New(createdGoParsedTimestamp)

	lastUpChangeGoParsedTimestamp, err := parseCustomTimestamp(osdDump.LastUpChange)
	if err != nil {
		return nil, err
	}
	lastUpChange := timestamppb.New(lastUpChangeGoParsedTimestamp)

	lastInChangeGoParsedTimestamp, err := parseCustomTimestamp(osdDump.LastInChange)
	if err != nil {
		return nil, err
	}
	lastInChange := timestamppb.New(lastInChangeGoParsedTimestamp)

	// Convert pools
	var osdDumpPools []*pb.OsdDumpPool
	for _, pool := range osdDump.Pools {
		createTimePbGoParsed, err := parseCustomTimestamp(pool.CreateTime)
		if err != nil {
			return nil, err
		}

		createTimePb := timestamppb.New(createTimePbGoParsed)

		osdDumpPools = append(osdDumpPools, &pb.OsdDumpPool{
			Pool:                              pool.Pool,
			PoolName:                          pool.PoolName,
			CreateTime:                        createTimePb,
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
		goParsedTime, err := parseCustomTimestamp(t)
		if err != nil {
			return nil, err
		}
		blocklistPb[ip] = timestamppb.New(goParsedTime)
	}

	var osdXInfo []*pb.OsdDumpOsdXInfo
	for _, osdX := range osdDump.OsdXinfo {
		downStamp, err := parseCustomTimestamp(osdX.DownStamp)
		if err != nil {
			return nil, err
		}
		lastPurgedSnapsScrub, err := parseCustomTimestamp(osdX.LastPurgedSnapsScrub)
		if err != nil {
			return nil, err
		}
		osdXInfo = append(osdXInfo, &pb.OsdDumpOsdXInfo{
			Osd:                  osdX.Osd,
			DownStamp:            timestamppb.New(downStamp),
			LaggyProbability:     osdX.LaggyProbability,
			LaggyInterval:        osdX.LaggyInterval,
			Features:             osdX.Features,
			OldWeight:            osdX.OldWeight,
			LastPurgedSnapsScrub: timestamppb.New(lastPurgedSnapsScrub),
			DeadEpoch:            osdX.DeadEpoch,
		})
	}

	response := pb.GetCephOsdDumpResponse{
		Epoch:                  osdDump.Epoch,
		Fsid:                   osdDump.Fsid,
		Modified:               modifiedTimestamp,
		Created:                createdTimestamp,
		LastUpChange:           lastUpChange,
		LastInChange:           lastInChange,
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

	return &response, nil
}

func parseCustomTimestamp(timestamp string) (time.Time, error) {
	const customTimeLayout = "2006-01-02T15:04:05.000000-0700"
	if timestamp == "0.000000" || timestamp == "" {
		// Return the zero time for Go, indicating an unset or invalid time.
		return time.Time{}, nil
	}
	return time.Parse(customTimeLayout, timestamp)
}
