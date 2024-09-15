//go:build windows
// +build windows

package collectorsmanager

import (
	sc "github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/collector"
)

func getCollectorsMap(l CollectorsList) []sc.StatisticCollector {
	return []sc.StatisticCollector{}
}
