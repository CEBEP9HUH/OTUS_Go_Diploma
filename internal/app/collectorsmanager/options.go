package collectorsmanager

type CollectorsList struct {
	EnableLoadAvg   bool `json:"enableLoadAvg"`
	EnableCPUUsage  bool `json:"enableCpuUsage"`
	EnableDiskUsage bool `json:"enableDiskUsage"`
	EnableNodeUsage bool `json:"enableNodeUsage"`
	EnableIOStat    bool `json:"enableIoStat"`
}
