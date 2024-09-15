package sysstatdeamon

import "github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/app/collectorsmanager"

type SysStatDaemonOpts struct {
	LogLevel   string                           `json:"logLevel"`
	Period     int                              `json:"periodSec"`
	Collectors collectorsmanager.CollectorsList `json:"collectors"`
	ServerPort uint
}
