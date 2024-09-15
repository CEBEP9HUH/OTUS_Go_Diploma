package grpcserver

import (
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/loggerwrapper"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/timedlist"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/pkg/api/grpc/statistic/broadcaster"
	"google.golang.org/grpc"
)

const (
	additionalStoreDuration = 5 * time.Second
)

var connIDCounter atomic.Int32

type grpcServer struct {
	broadcaster.UnimplementedSysStatBroadcasterServer
	storeDurations []time.Duration
	timeShift      time.Duration
	storage        timedlist.TimedList[statistic.Snapshot]
	logger         loggerwrapper.Logger
	m              sync.Mutex
}

func makeServer(storage timedlist.TimedList[statistic.Snapshot],
	logger loggerwrapper.Logger, timeShift time.Duration,
) broadcaster.SysStatBroadcasterServer {
	return &grpcServer{
		storage:   storage,
		logger:    logger,
		timeShift: timeShift,
	}
}

func (s *grpcServer) Subscribe(params *broadcaster.StatParams,
	stream grpc.ServerStreamingServer[broadcaster.SysStat],
) error {
	connID := connIDCounter.Add(1)
	s.logger.Info("Receive new connection (%d)", connID)
	defer s.logger.Info("Close connection %d", connID)
	h := s.initConnHandler(connID, params)
	defer s.deinitConnHandler(h)
	return h.handle(stream)
}

func (s *grpcServer) initConnHandler(id int32, params *broadcaster.StatParams) connHandler {
	l, err := loggerwrapper.NewStdLogger(fmt.Sprintf("Conn handler #%d", id), s.logger.Level())
	if err != nil {
		s.logger.Fatal(err.Error())
	}
	h := connHandler{
		storage:          s.storage,
		loger:            l,
		collectingPeriod: params.CollectingPeriod.AsDuration(),
		sendingPeriod:    params.SendingPeriod.AsDuration(),
		timeShift:        s.timeShift,
	}
	s.m.Lock()
	defer s.m.Unlock()
	s.storeDurations = append(s.storeDurations, h.collectingPeriod)
	s.storage.SetStoreDuration(slices.Max(s.storeDurations) + additionalStoreDuration)
	return h
}

func (s *grpcServer) deinitConnHandler(h connHandler) {
	s.m.Lock()
	defer s.m.Unlock()
	pos := slices.Index(s.storeDurations, h.collectingPeriod)
	s.storeDurations = slices.Delete(s.storeDurations, pos, pos+1)
	if len(s.storeDurations) == 0 {
		s.storage.SetStoreDuration(additionalStoreDuration)
		return
	}
	s.storage.SetStoreDuration(slices.Max(s.storeDurations) + additionalStoreDuration)
}
