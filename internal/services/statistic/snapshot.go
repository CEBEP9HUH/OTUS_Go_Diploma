package statistic

type Snapshot struct {
	CPUUsage  CPUUsage
	DiskUsage DiskUsage
	IOStat    IOStat
	NodeUsage NodeUsage
	LoadAvg   LoadAvg
}

func (snap *Snapshot) Add(stat Statistic) {
	switch statData := stat.(type) {
	case CPUUsage:
		snap.CPUUsage = statData
	case DiskUsage:
		snap.DiskUsage = statData
	case IOStat:
		snap.IOStat = statData
	case NodeUsage:
		snap.NodeUsage = statData
	case LoadAvg:
		snap.LoadAvg = statData
	}
}
