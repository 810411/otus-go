package grpc

import (
	"time"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/api/grpcpb"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func mapEventRepoToPb(v repository.Event) *grpcpb.Event {
	return &grpcpb.Event{
		Id:       uint64(v.ID),
		Title:    v.Title,
		Datetime: &timestamp.Timestamp{Seconds: v.Datetime.Unix(), Nanos: int32(v.Datetime.Nanosecond())},
		Duration: &duration.Duration{Seconds: int64(v.Duration.Seconds())},
		OwnerId:  v.OwnerID,
	}
}

func mapEventPbToRepo(v *grpcpb.Event) repository.Event {
	return repository.Event{
		ID:          repository.EventID(v.Id),
		Title:       v.Title,
		Datetime:    time.Unix(v.Datetime.Seconds, int64(v.Datetime.Nanos)),
		Duration:    time.Duration(v.Duration.Seconds),
		Description: v.Description,
		OwnerID:     v.OwnerId,
	}
}
