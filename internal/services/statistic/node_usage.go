package statistic

import (
	"strconv"
	"strings"
)

type FSNodeInfo struct {
	FS    string
	Usage float32
}

type NodeUsage struct {
	NodeUsage map[string]FSNodeInfo
}

func (NodeUsage) isStatistic() {}

func (sd NodeUsage) String() string {
	var builder strings.Builder
	for mount, info := range sd.NodeUsage {
		builder.WriteString(mount)
		builder.WriteString(": ")
		builder.WriteString(info.FS)
		builder.WriteString(" - ")
		builder.WriteString(strconv.FormatFloat(float64(info.Usage), 'f', 2, 32))
		builder.WriteString("%;")
	}
	return builder.String()
}
