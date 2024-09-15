package statistic

import (
	"strconv"
	"strings"
)

type FSDiskInfo struct {
	FS    string
	Usage float32
}

type DiskUsage struct {
	BlockUsage map[string]FSDiskInfo
}

func (DiskUsage) isStatistic() {}

func (sd DiskUsage) String() string {
	var builder strings.Builder
	for mount, info := range sd.BlockUsage {
		builder.WriteString(mount)
		builder.WriteString(": ")
		builder.WriteString(info.FS)
		builder.WriteString(" - ")
		builder.WriteString(strconv.FormatFloat(float64(info.Usage), 'f', 2, 32))
		builder.WriteString("%;")
	}
	return builder.String()
}
