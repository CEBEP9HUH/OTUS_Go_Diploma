package statistic

import (
	"strconv"
	"strings"
)

type DevIOStat struct {
	TPS      float32
	ReadKbs  float32
	WriteKbs float32
}

type IOStat struct {
	Data map[string]DevIOStat
}

func (IOStat) isStatistic() {}

func (sd IOStat) String() string {
	var builder strings.Builder
	for dev, info := range sd.Data {
		builder.WriteString(dev)
		builder.WriteString(": ")
		builder.WriteString(strconv.FormatFloat(float64(info.TPS), 'f', 2, 32))
		builder.WriteString("; ")
		builder.WriteString(strconv.FormatFloat(float64(info.ReadKbs), 'f', 2, 32))
		builder.WriteString("; ")
		builder.WriteString(strconv.FormatFloat(float64(info.WriteKbs), 'f', 2, 32))
		builder.WriteString(";")
	}
	return builder.String()
}
