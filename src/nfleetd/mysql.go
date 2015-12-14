package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Provides a set of configuration variables that describe how to interact with a SQL database.
type SQLConf struct {
	driver  string
	openStr string
	table   string
	insertTable string
}

const (
	driver = "mysql"
	openStr = "nfleet:dpsvmflxm@tcp(www.motrexlab.net:3306)/nfleet"
	table = "geofence"
	insertTable = "inout_geofence"
	latCol = "lat"
	lngCol = "lng"
)

// Returns a SQLConf based on the $DB environment variable
func GetSqlConf() *SQLConf {
	return &SQLConf{driver: driver, openStr: openStr, table: table, insertTable : insertTable}
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
