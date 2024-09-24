package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
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

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	cfgPath := filepath.Join(exPath, "config.json")
	err = config.InitConfig(cfgPath)
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
