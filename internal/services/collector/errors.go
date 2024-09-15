package collector

import "errors"

// Общие ошибки сборщиков статистики.
var (
	ErrCPUInfoNotFound       = errors.New("can't find cpu info")
	ErrUserModeNotFound      = errors.New("can't find user mode info")
	ErrSystemModeNotFound    = errors.New("can't find system mode info")
	ErrIdleNotFound          = errors.New("can't find idle info")
	ErrDiskUsageInfoNotFound = errors.New("can't find disk usage info")
	ErrNodeUsageInfoNotFound = errors.New("can't find node usage info")
	ErrLoadAvgInfoNotFound   = errors.New("can't find load avereage info")
	ErrAvgMin1NotFound       = errors.New("can't find 1 minute load avereage info")
	ErrAvgMin5NotFound       = errors.New("can't find 5 minute load avereage info")
	ErrAvgMin15NotFound      = errors.New("can't find 15 minute load avereage info")
	ErrTPSInfoNotFound       = errors.New("can't find transfers per second info")
	ErrReadKbsInfoNotFound   = errors.New("can't find kilobytes read per second info")
	ErrWriteKbsInfoNotFound  = errors.New("can't find kilobytes write per second info")
)
