package statistic

import "fmt"

type CPUUsage struct {
	UserMode   float32
	SystemMode float32
	Idle       float32
}

func (CPUUsage) isStatistic() {}

func (sd CPUUsage) String() string {
	return fmt.Sprintf("%.1f us, %.1f sy, %.1f id", sd.UserMode, sd.SystemMode, sd.Idle)
}
