package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/app/sysstatdeamon"
	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/config"
)

var signalsList = []os.Signal{os.Interrupt, syscall.SIGHUP, syscall.SIGTERM}

func main() {
	defer func() {
		if recoverError := recover(); recoverError != nil {
			fmt.Printf("Got panic, exiting: %s\n", recoverError)
		}
	}()
	ctx, cancel := signal.NotifyContext(context.Background(), signalsList...)
	defer cancel()

	err := config.InitConfig("config.json")
	if err != nil {
		panic(err)
	}

	cm, err := sysstatdeamon.MakeSysStatDaemon(config.GetConfig().SysStatDaemonOpts)
	if err != nil {
		panic(err)
	}
	if err := cm.Run(ctx); err != nil {
		panic(err)
	}
}
