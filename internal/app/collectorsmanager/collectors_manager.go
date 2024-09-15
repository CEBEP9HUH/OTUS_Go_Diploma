package collectorsmanager

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/collector"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/defaults"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/loggerwrapper"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
)

const (
	minimumCollectInterval = defaults.MinStatInterval
)

var (
	ErrEmptyCollectorsList  = errors.New("all statistic collectors are disabled")
	ErrEmptyStatisticOutput = errors.New("statistic output can't be nil")
)

type collectorsManager struct {
	statisticOutput chan<- statistic.Statistic
	period          time.Duration
	collectors      []collector.StatisticCollector
	logger          loggerwrapper.Logger
}

type CollectorsManager interface {
	Run(ctx context.Context) error
}

func MakeCollectorsManager(statisticOutput chan<- statistic.Statistic, collectInterval time.Duration,
	l loggerwrapper.Logger, list CollectorsList,
) (CollectorsManager, error) {
	collectors := getCollectorsMap(list)
	if len(collectors) == 0 {
		return nil, ErrEmptyCollectorsList
	}
	if collectInterval < minimumCollectInterval {
		l.Warn("Collect interval can't be less than %q. Set interval to minimum value.", minimumCollectInterval)
		collectInterval = minimumCollectInterval
	}

	return &collectorsManager{
		statisticOutput: statisticOutput,
		collectors:      collectors,
		period:          collectInterval,
		logger:          l,
	}, nil
}

func (m *collectorsManager) Run(ctx context.Context) error {
	m.logger.Info("Run collectors manager")
	defer m.logger.Info("Stop collectors manager")
	collect := func(ctx context.Context, wg *sync.WaitGroup, collector collector.StatisticCollector) {
		defer wg.Done()
		for {
			localCtx, cancel := context.WithCancel(ctx)
			go m.collect(localCtx, collector)
			select {
			case <-ctx.Done():
				cancel()
				return
			case <-time.After(m.period):
				cancel()
			}
		}
	}
	var wg sync.WaitGroup
	for _, c := range m.collectors {
		wg.Add(1)
		go collect(ctx, &wg, c)
	}
	wg.Wait()
	return nil
}

// private section

func (m *collectorsManager) collect(ctx context.Context, collector collector.StatisticCollector) {
	stat, err := collector.Collect(ctx)
	if err != nil {
		m.logger.Error("Unable to collect data: %s", err)
		return
	}
	if stat == nil {
		m.logger.Warn("Stat is nil. Skip")
		return
	}
	select {
	case <-ctx.Done():
		return
	case m.statisticOutput <- stat:
	}
}
