package main

import (
	"database/sql"
	"fmt"
)

// A Mapper that uses Standard SQL Syntax to perform mapping functions and queries
type SQLMapper struct {
	conf    *SQLConf
	sqlConn *sql.DB
}

// Returns a pointer to the SQLMapper's SQL Database Connection.
func (s *SQLMapper) SqlDbConn() *sql.DB {
	return s.sqlConn
}

func (s *SQLMapper) GetDeviceListFromDatabase() map[string]int{

	query := fmt.Sprintf("SELECT imei FROM %v", s.conf.devicetable)

	rows, err := s.sqlConn.Query(query)
	if err != nil{
		panic(err)
	}

	deviceList := make(map[string]int)

	var imei string

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&imei)
		if err != nil {
			log.Fatal(err)
		}
		deviceList[imei] = 0;
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return deviceList
}


func (s *SQLMapper) GetGeofenceListFromDatabase() []Geofence{

	query := fmt.Sprintf("SELECT id, shape, vertices, radius FROM %v where active=1", s.conf.geofencetable)

	rows, err := s.sqlConn.Query(query)
	if err != nil{
		panic(err)
	}

	geofenceList := make([]Geofence, 0)
	geofence := new(Geofence);

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&geofence.id, &geofence.shape, &geofence.vertices, &geofence.radius)
		if err != nil {
			log.Fatal(err)
		}
		geofence.state = make(map[string]int)
		geofenceList = append(geofenceList, *geofence)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return geofenceList
}

func (s *SQLMapper) InsertResult(resultList []Geofence) error{

	query := fmt.Sprintf("INSERT INTO %s(geofence_id, imei, time, flag) VALUES(?, ?, now(), ?)", s.conf.inouttable)

	stmt, err := s.sqlConn.Prepare(query)
	if err != nil{
		log.Fatal(err)
	}

	for _, result := range resultList{
		for imei, flag := range result.state{
			_, err := stmt.Exec(result.id, imei, flag)
			if err != nil{
				log.Fatal(err)
			}
		}
	}

	return err
}

func (s *SQLMapper) GetGeofenceState(geofenceList []Geofence){

	query := fmt.Sprintf("SELECT imei, flag FROM %v  where geofence_id = ? group by geofence_id, imei order by time desc;", s.conf.inouttable)

	stmt, err := s.sqlConn.Prepare(query)
	if err != nil{
		log.Fatal(err)
	}

	for _, geofence := range geofenceList{

		var imei string
		var flag int


		rows, err := stmt.Query(geofence.id)

		if err != nil{
			log.Fatal(err)
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&imei, &flag)
			if err != nil {
				log.Fatal(err)
			}
			geofence.state[imei] = flag
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *SQLMapper) InsertDeivceInfoToMysql(imei string){
	query := fmt.Sprintf("INSERT INTO %s(imei) VALUES(?)", s.conf.devicetable)

	stmt, err := s.sqlConn.Prepare(query)
	if err != nil{
		log.Fatal(err)
	}

	_, err = stmt.Exec(imei)
	if err != nil{
		log.Fatal(err)
	}
}




