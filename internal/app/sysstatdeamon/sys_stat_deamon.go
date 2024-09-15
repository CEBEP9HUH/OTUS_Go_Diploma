package sysstatdeamon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/app/broadcaster"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/app/broadcaster/grpcserver"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/app/collectorsmanager"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/app/statisticmanager"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/loggerwrapper"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/timedlist"
)

type deinitCalls struct {
	calls []func()
}

type servers struct {
	broadcaster broadcaster.Broadcaster
}

type communications struct {
	statistic chan statistic.Statistic
}

type storages struct {
	storage timedlist.TimedList[statistic.Snapshot]
}

type managers struct {
	statisticManager  statisticmanager.StatisticManager
	collectorsManager collectorsmanager.CollectorsManager
}

type sysStatDaemon struct {
	SysStatDaemonOpts

	logger         loggerwrapper.Logger
	communications communications
	storages       storages
	managers       managers
	servers        servers
	deinitCalls    deinitCalls
}

type SysStatDaemon interface {
	Run(ctx context.Context) error
}

func MakeSysStatDaemon(opts SysStatDaemonOpts) (SysStatDaemon, error) {
	const (
		errTemplate = "statistic daemon creation: %w"
	)
	l, err := loggerwrapper.NewStdLogger("statistic daemon", opts.LogLevel)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	ret := &sysStatDaemon{
		logger:            l,
		SysStatDaemonOpts: opts,
	}
	if err := ret.init(); err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}
	return ret, nil
}

func (s *sysStatDaemon) Run(ctx context.Context) error {
	s.logger.Info("Run statistic daemon")
	defer s.logger.Info("Stop statistic daemon")
	defer func() {
		for _, deinit := range s.deinitCalls.calls {
			deinit()
		}
	}()
	var wg sync.WaitGroup
	localCtx, cancel := context.WithCancel(ctx)
	s.deinitCalls.calls = append(s.deinitCalls.calls, cancel)

	runManager := func(name string, manager interface{ Run(context.Context) error }) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := manager.Run(localCtx); err != nil {
				s.logger.Error("%s finished with an error: %s", name, err)
				cancel()
			}
		}()
	}

	runManager("Collectors manager", s.managers.collectorsManager)
	runManager("Statistic manager", s.managers.statisticManager)
	runManager("Broadcaster", s.servers.broadcaster)

	wg.Wait()
	return nil
}

// private section

func (s *sysStatDaemon) init() error {
	const (
		errTemplate = "statistic daemon init: %w"
	)
	s.initCommunications()

	period := time.Duration(s.Period) * time.Second
	s.initStorage(period)
	if err := s.initCollectorsManager(period); err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	if err := s.initStatisticManager(period); err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	if err := s.initServer(period); err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	return nil
}

func (s *sysStatDaemon) initCommunications() {
	s.communications.statistic = make(chan statistic.Statistic)
	s.deinitCalls.calls = append(s.deinitCalls.calls,
		func() {
			close(s.communications.statistic)
		},
	)
}

func (s *sysStatDaemon) initStorage(period time.Duration) {
	st, cancel := timedlist.MakeTimedList[statistic.Snapshot](period * 5)
	s.storages.storage = st
	s.deinitCalls.calls = append(s.deinitCalls.calls, cancel)
}

func (s *sysStatDaemon) initCollectorsManager(period time.Duration) error {
	const (
		errTemplate = "collectors manager init: %w"
	)
	l, err := loggerwrapper.NewStdLogger("collectors manager", s.LogLevel)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	m, err := collectorsmanager.MakeCollectorsManager(s.communications.statistic, period, l, s.Collectors)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	s.managers.collectorsManager = m
	return nil
}

func (s *sysStatDaemon) initStatisticManager(period time.Duration) error {
	const (
		errTemplate = "statistic manager init: %w"
	)
	l, err := loggerwrapper.NewStdLogger("statistic manager", s.LogLevel)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	m, err := statisticmanager.MakeStatisticManager(s.communications.statistic, s.storages.storage, period, l)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	s.managers.statisticManager = m
	return nil
}

func (s *sysStatDaemon) initServer(period time.Duration) error {
	const (
		errTemplate = "server init: %w"
	)
	l, err := loggerwrapper.NewStdLogger("statistic broadcaster", s.LogLevel)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	m, err := grpcserver.MakeGRPCBroadcaster(s.storages.storage, l, s.ServerPort, period)
	if err != nil {
		return fmt.Errorf(errTemplate, err)
	}
	s.servers.broadcaster = m
	return nil
}
