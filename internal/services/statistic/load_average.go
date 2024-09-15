package statistic

import "fmt"

type LoadAvg struct {
	Min1  float32
	Min5  float32
	Min15 float32
}

func (LoadAvg) isStatistic() {}

func (sd LoadAvg) String() string {
	return fmt.Sprintf("%.1f, %.1f, %.1f", sd.Min1, sd.Min5, sd.Min15)
}
