package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/app/broadcaster"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/loggerwrapper"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/timedlist"
	gbroadcaster "github.com/CEBEP9HUH/OTUS_Go_Diploma/pkg/api/grpc/statistic/broadcaster"
	"google.golang.org/grpc"
)

type grpcBroadcaster struct {
	server gbroadcaster.SysStatBroadcasterServer
	logger loggerwrapper.Logger
	port   uint
}

func MakeGRPCBroadcaster(storage timedlist.TimedList[statistic.Snapshot],
	logger loggerwrapper.Logger, port uint, timeShift time.Duration,
) (broadcaster.Broadcaster, error) {
	l, err := loggerwrapper.NewStdLogger("Server", logger.Level())
	if err != nil {
		return nil, fmt.Errorf("unable create server: %w", err)
	}
	return &grpcBroadcaster{
		server: makeServer(storage, l, timeShift),
		port:   port,
		logger: logger,
	}, nil
}

func (b *grpcBroadcaster) Run(ctx context.Context) error {
	b.logger.Info("Run broadcaster")
	defer b.logger.Info("Stop broadcaster")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", b.port))
	if err != nil {
		return fmt.Errorf("unable to start broadcaster: %w", err)
	}
	grpcServer := grpc.NewServer()
	gbroadcaster.RegisterSysStatBroadcasterServer(grpcServer, b.server)
	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()
	return grpcServer.Serve(lis)
}
