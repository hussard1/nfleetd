package main

import (
	"flag"
	"runtime"
)

func commandLine() (configfile *string, logfile *string, loglevel *string) {
	configfile = flag.String("config", "", "specify the nfleetd config file")
	logfile = flag.String("logfile", "", "file where the nfleetd log will be stored (default \"console\")")
	loglevel = flag.String("loglevel", "info", "log level")

	flag.Parse()

	return
}

func bootup(worker *Worker, devices []Device) {
	log.Info("Starting nfleetd server...")

	server := new(Server)

	for _, device := range devices {
		if device.enabled == true {
			go func(w Worker, d Device) {
				server.Bind(w, d)
			}(*worker, device)
		}
	}
}

func main() {
	// prevent to stop main goroutine
	done := make(chan struct {})
	defer close(done)

	// use all cores
	numcpus := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpus)

	configfile, logfile, loglevel := commandLine()

	// Initialize logger
	InitLogger(logfile, loglevel)

	// Load configuration file
	conf := new(Configuration)
	conf.Load(configfile)

	bootup(conf.getWorker(), conf.GetDevices())

	for _ = range done {
		select {
		case <-done:
			return
		}
	}
}