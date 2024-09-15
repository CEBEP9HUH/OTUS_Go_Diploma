package config

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"os"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/app/sysstatdeamon"
)

type Config struct {
	SysStatDaemonOpts sysstatdeamon.SysStatDaemonOpts `json:"sysStatDaemon"`
}

var config Config

func InitConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	var data []byte
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		data = append(data, buf[:n]...)
	}
	if err := InitConfigFromBytes(data); err != nil {
		return err
	}
	flag.Parse()
	config.SysStatDaemonOpts.ServerPort = serverPort
	return nil
}

func InitConfigFromBytes(data []byte) error {
	return json.Unmarshal(data, &config)
}

func GetConfig() Config {
	return config
}
