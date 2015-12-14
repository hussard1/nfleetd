package main

import (
	"flag"
	"runtime"
	"gopkg.in/mgo.v2"
)

func commandLine() (configfile *string, logfile *string, loglevel *string) {
	configfile = flag.String("config", "", "specify the nfleetd config file")
	logfile = flag.String("logfile", "", "file where the nfleetd log will be stored (default \"console\")")
	loglevel = flag.String("loglevel", "info", "log level")

	flag.Parse()

	return
}

func bootup(worker *Worker, devices []Device, geofenceList []Geofence, collection *mgo.Session, mysqlSession *SQLMapper) {
	log.Info("Starting nfleetd server...")

	server := new(Server)

	for _, device := range devices {
		if device.enabled == true {
			go func(w Worker, d Device, g []Geofence, s *mgo.Session, ms *SQLMapper){
				server.Bind(w, d, g, s, ms)
			}(*worker, device, geofenceList, collection, mysqlSession)
		}
	}
}

func InitGeofence(mysqlSession *SQLMapper) ([]Geofence){

	// Get Geofence List
	geofenceList := mysqlSession.GetGeofenceFromDatabase()

	// Get Inoutstate
	mysqlSession.GetGeofenceState(geofenceList)


	// initialize Inoutstate
	return geofenceList
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
	// Initialize mongoDB
	session := InitMongoDB()

	// Initialize mysql
	mysqlSession := GetSQLMapper()

	// Initialize geofence
	geofenceList := InitGeofence(mysqlSession)


	bootup(conf.getWorker(), conf.GetDevices(), geofenceList, session, mysqlSession)

	for _ = range done {
		select {
		case <-done:
			return
		}
	}
}