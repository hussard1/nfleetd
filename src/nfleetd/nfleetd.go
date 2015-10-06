package main

import (
	"flag"
)

func commandLine() (configfile *string, logfile *string, loglevel *string) {
	configfile = flag.String("config", "", "specify the nfleetd config file")
	logfile = flag.String("logfile", "", "file where the nfleetd log will be stored (default \"console\")")
	loglevel = flag.String("loglevel", "info", "log level")

	flag.Parse()

	return
}

func bootup() {
	log.Info("Starting nfleetd server...")
}

func main() {
	configfile, logfile, loglevel := commandLine()

	// Initialize logger
	InitLogger(logfile, loglevel)

	// Load configuration file
	conf := new(Configuration)
	conf.Load(configfile)

	bootup()
}