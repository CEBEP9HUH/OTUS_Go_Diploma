package main

import (
	"flag"
	"time"
)

const (
	CPU = iota
	IO
	Load
	DiskUsage
	NodeUsage
)

var (
	address         string
	collectDuration time.Duration
	sendDuration    time.Duration
	statType        uint64
)

func init() {
	const (
		defaultAddr            = "localhost:5678"
		defaultCollectDuration = time.Second * 2
		defaultSendDuration    = time.Second * 10
		defaultStatType        = 0
	)
	flag.StringVar(&address, "addr", defaultAddr, "statistic server address:port")
	flag.DurationVar(&collectDuration, "collect", defaultCollectDuration, "statistic collect duration")
	flag.DurationVar(&sendDuration, "interval", defaultSendDuration, "statistic sending interval")
	flag.Uint64Var(&statType, "print", defaultStatType,
		"statistic to print\n\t0 - cpu\n\t1 - io\n\t2 - load\n\t3 - disk\n\t4 - node")
}
