package api

import "google.golang.org/protobuf/types/known/timestamppb"

func tsToPb(in *int) *timestamppb.Timestamp {
	if in == nil {
		return nil
	}
	return &timestamppb.Timestamp{Seconds: int64(*in)}
}
