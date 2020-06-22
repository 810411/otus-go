package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc/peer"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/api/grpcpb"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/api"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
	"google.golang.org/grpc"
)

type Service struct {
	ctx      context.Context
	repo     repository.EventsRepo
	settings api.Settings
	server   *grpc.Server
}

func New(c context.Context, s api.Settings, r repository.EventsRepo) *Service {
	service := &Service{
		ctx:      c,
		repo:     r,
		settings: s,
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logInterceptor))
	grpcpb.RegisterEventsServer(grpcServer, service)
	service.server = grpcServer

	return service
}

func (s *Service) ListenAndServe() error {
	lsn, err := net.Listen("tcp", net.JoinHostPort(s.settings.Host, s.settings.Port))
	if err != nil {
		return err
	}

	if err := s.server.Serve(lsn); err != nil {
		return err
	}

	return nil
}

func (s *Service) GracefulStop() {
	s.server.GracefulStop()
}

func logInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	h, err := handler(ctx, req)
	latency := time.Since(startTime)

	var remoteAddr string
	if p, ok := peer.FromContext(ctx); ok {
		remoteAddr = p.Addr.String()
	}

	logger.Logger.Info(fmt.Sprintf("%s [%s] %s %s %s \"%s\"",
		remoteAddr,
		time.Now().Format("2006-01-02 15:04:05 -0700"),
		info.FullMethod,
		req,
		latency,
		"GRPC-client",
	))

	return h, err
}
