package config

import "flag"

var serverPort uint

func init() {
	const (
		defaultPort uint = 5678
	)
	flag.UintVar(&serverPort, "port", defaultPort, "port for listening")
}
