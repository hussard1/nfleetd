package main

import (
	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()
	log.Info("Starting up nfleetd server...")
}