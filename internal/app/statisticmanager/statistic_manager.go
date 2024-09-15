package statisticmanager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/defaults"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/loggerwrapper"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/timedlist"
)

const (
	minimumStoreInterval = defaults.MinStatInterval
)

var (
	ErrEmptyStatisticInput = errors.New("statistic input can't be nil")
	ErrStorageUnset        = errors.New("storage can't be nil")
)

type statisticManager struct {
	statisticInput <-chan statistic.Statistic
	storeInterval  time.Duration
	storage        timedlist.TimedList[statistic.Snapshot]
	logger         loggerwrapper.Logger
}

type StatisticManager interface {
	Run(ctx context.Context) error
}

func MakeStatisticManager(statisticInput <-chan statistic.Statistic, storage timedlist.TimedList[statistic.Snapshot],
	storeInterval time.Duration, l loggerwrapper.Logger,
) (StatisticManager, error) {
	if statisticInput == nil {
		return nil, ErrEmptyStatisticInput
	}
	if storage == nil {
		return nil, ErrStorageUnset
	}
	if storeInterval < minimumStoreInterval {
		l.Warn("Store interval can't be less than %q. Set interval to minimum value.", minimumStoreInterval)
		storeInterval = minimumStoreInterval
	}
	return &statisticManager{
		statisticInput: statisticInput,
		storeInterval:  storeInterval,
		storage:        storage,
		logger:         l,
	}, nil
}

func (m *statisticManager) Run(ctx context.Context) error {
	m.logger.Info("Run statistic manager")
	defer m.logger.Info("Stop statistic manager")
	var snap statistic.Snapshot
	ticker := time.NewTicker(m.storeInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case stat, ok := <-m.statisticInput:
			if !ok {
				return fmt.Errorf("statistic channel was closed")
			}
			snap.Add(stat)
		case <-ticker.C:
			m.storage.Add(snap, time.Now())
			snap = statistic.Snapshot{}
		}
	}
}
