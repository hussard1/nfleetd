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

func bootup(worker *Worker, devices []Device, deviceList map[string]int, geofenceList []Geofence, mongoSesson *MongoSession, mysqlSession *SQLMapper) {
	log.Info("Starting nfleetd server...")

	server := new(Server)

	for _, device := range devices {
		if device.enabled == true {
			go func(w Worker, d Device, dl map[string]int,  g []Geofence, ms *MongoSession, s *SQLMapper){
				server.Bind(w, d, dl, g, ms, s)
			}(*worker, device, deviceList, geofenceList, mongoSesson, mysqlSession)
		}
	}
}

func InitDeviceList(mysqlSession *SQLMapper) map[string]int{

	deviceList := mysqlSession.GetDeviceListFromDatabase();

	return deviceList
}

func InitGeofence(mysqlSession *SQLMapper) ([]Geofence){

	// Get Geofence List
	geofenceList := mysqlSession.GetGeofenceListFromDatabase()

	// Get Inoutstate
	mysqlSession.GetGeofenceState(geofenceList)

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
	mongoSession := InitMongoDB()

	// Initialize mysql
	mysqlSession := GetSQLMapper()

	// Initailize device list
	deviceList := InitDeviceList(mysqlSession)

	// Initialize geofence
	geofenceList := InitGeofence(mysqlSession)


	bootup(conf.getWorker(), conf.GetDevices(), deviceList, geofenceList, mongoSession, mysqlSession)

	for _ = range done {
		select {
		case <-done:
			return
		}
	}
}