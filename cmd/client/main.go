package main

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/pkg/api/grpc/statistic/broadcaster"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := broadcaster.NewSysStatBroadcasterClient(conn)
	p := &broadcaster.StatParams{
		SendingPeriod:    durationpb.New(sendDuration),
		CollectingPeriod: durationpb.New(collectDuration),
	}
	stream, err := client.Subscribe(context.Background(), p)
	if err != nil {
		panic(err)
	}

	toOutput := make(chan *broadcaster.SysStat)
	defer close(toOutput)
	switch statType {
	case CPU:
		go printCPUusage(toOutput)
	case IO:
		go printIOStat(toOutput)
	case Load:
		go printLoadAvg(toOutput)
	case DiskUsage:
		go printDiskUsage(toOutput)
	case NodeUsage:
		go printNodeUsage(toOutput)
	default:
		panic("unknown type to print")
	}
	for {
		v, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		toOutput <- v
	}
}
