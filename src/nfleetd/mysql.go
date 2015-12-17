package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Provides a set of configuration variables that describe how to interact with a SQL database.
type SQLConf struct {
	driver  string
	openStr string
	devicetable string
	geofencetable   string
	inouttable string
}

const (
	driver = "mysql"
	openStr = "nfleet:dpsvmflxm@tcp(nfleet.c3li2b1mmcbj.ap-southeast-1.rds.amazonaws.com:3306)/nfleet"
	devicetable = "device"
	geofencetable = "geofence"
	inouttable = "inout_geofence"
	latCol = "lat"
	lngCol = "lng"
)

// Returns a SQLConf based on the $DB environment variable
func GetSqlConf() *SQLConf {
	return &SQLConf{driver: driver, openStr: openStr, devicetable: devicetable, geofencetable: geofencetable, inouttable : inouttable}
}

func GetSQLMapper() (*SQLMapper) {
	sqlConf := GetSqlConf()

	s := &SQLMapper{conf: sqlConf}

	db, err := sql.Open(s.conf.driver, s.conf.openStr)
	if err != nil {
		log.Panic("Create Mysql Session: %s\n", err)
	}

	s.sqlConn = db
	return s
}
