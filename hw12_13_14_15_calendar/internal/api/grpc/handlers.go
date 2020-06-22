package grpc

import (
	"context"
	"time"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/api/grpcpb"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
)

func (s *Service) Create(ctx context.Context, req *grpcpb.CreateRequest) (*grpcpb.CreateResponse, error) {
	event := mapEventPbToRepo(req.Event)

	event, err := s.repo.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	return &grpcpb.CreateResponse{
		Event: mapEventRepoToPb(event),
	}, nil
}

func (s *Service) Update(ctx context.Context, req *grpcpb.UpdateRequest) (*grpcpb.UpdateResponse, error) {
	id := repository.EventID(req.Id)
	event := mapEventPbToRepo(req.Event)
	event.ID = id

	event, err := s.repo.Update(ctx, event)
	if err != nil {
		return nil, err
	}

	return &grpcpb.UpdateResponse{
		Event: mapEventRepoToPb(event),
	}, nil
}

func (s *Service) Delete(ctx context.Context, req *grpcpb.DeleteRequest) (*grpcpb.DeleteResponse, error) {
	id := repository.EventID(req.Id)

	id, err := s.repo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return &grpcpb.DeleteResponse{
		Id: uint64(id),
	}, nil
}

func (s *Service) ListOfDate(ctx context.Context, req *grpcpb.ListOfRequest) (*grpcpb.ListOfResponse, error) {
	return listOf(ctx, req, s.repo.ListOfDay)
}

func (s *Service) ListOfWeek(ctx context.Context, req *grpcpb.ListOfRequest) (*grpcpb.ListOfResponse, error) {
	return listOf(ctx, req, s.repo.ListOfWeek)
}

func (s *Service) ListOfMonth(ctx context.Context, req *grpcpb.ListOfRequest) (*grpcpb.ListOfResponse, error) {
	return listOf(ctx, req, s.repo.ListOfMonth)
}

func listOf(ctx context.Context, req *grpcpb.ListOfRequest, fn func(ctx context.Context, from time.Time) ([]repository.Event, error)) (*grpcpb.ListOfResponse, error) {
	from := time.Unix(req.From.Seconds, int64(req.From.Nanos))
	eventsRepo, err := fn(ctx, from)
	if err != nil {
		return nil, err
	}

	events := make([]*grpcpb.Event, len(eventsRepo))
	for i, v := range eventsRepo {
		event := mapEventRepoToPb(v)
		events[i] = event
	}

	return &grpcpb.ListOfResponse{Event: events}, nil
}
