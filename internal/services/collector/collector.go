package collector

import (
	"context"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
)

type StatisticCollector interface {
	Collect(ctx context.Context) (statistic.Statistic, error)
}
