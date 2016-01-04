package main

import (
	"encoding/json"
	"bytes"
	"rule"
	"database/sql"
)


type Geofence struct {
	id       int
	shape    int
	vertices string
	radius   sql.NullFloat64
	state 	 map[string]int
}


//func checkInOutGeofence(geofenceList []Geofence, msgList []rule.Message) []Geofence{
func CheckInOutGeofence(msgList []rule.Message, geofenceList []Geofence, mysqlSession *SQLMapper){

//	var IMEI string
//	for _, msg := range msgList {
//		IMEI = msg.IMEI
//	}
	//select geofence state
//	mysqlSession.GetGeofenceState(geofenceList)

	resultList := Checkpoint(geofenceList, msgList)
	if resultList != nil && len(resultList) != 0{
		_ = mysqlSession.InsertResult(resultList)
		//update database
	}
}

func Checkpoint(geofenceList []Geofence, msgList []rule.Message) []Geofence{

	resultList := make([]Geofence, 0)

	for _, msg := range msgList{
		point := NewPoint(msg.Latitude, msg.Longitude)
		for _, geofence := range geofenceList{
			if geofence.shape == 1{
				if PointsWithinCircle(geofence, point){
					if CheckInoutState(geofence, msg.IMEI){
						resultList = append(resultList, geofence)
						geofence.state[msg.IMEI] = 2
					}
				}else{
					if geofence.state[msg.IMEI] == 2{
						geofence.state[msg.IMEI] = 1
						resultList = append(resultList, geofence)
					}
				}
			}else if geofence.shape == 2{
				if PointsWithinPolygon(geofence, point) {
					if CheckInoutState(geofence, msg.IMEI){
						resultList = append(resultList, geofence)
						geofence.state[msg.IMEI] = 2
					}
				}else{
					if geofence.state[msg.IMEI] ==2{
						geofence.state[msg.IMEI] = 1
						resultList = append(resultList, geofence)
					}
				}
			}
		}
	}

	return resultList
}


func CheckInoutState(geofence Geofence, imei string) bool{

	flag, _ := geofence.state[imei]

	if flag == 1 || flag == 0{
		 return true
	}else if flag ==2 {
		return false
	}else{
		return false
	}
}


func PointsWithinCircle(geofence Geofence, point *Point) bool{

	var geoPoint *Point
	geoPoint = new(Point)
	geoPoint.UnmarshalJSON(geofence.vertices)
	distance := point.GetCircleDistance(geoPoint)
	if distance <= geofence.radius.Float64{
		return true;
	}else{
		return false;
	}
}

func PointsWithinPolygon(geofence Geofence, point *Point) bool{

	var geoPoint *Point
	geoPoint = new(Point)
	values := geoPoint.UnmarshalJSONArray(geofence.vertices)
	polygon := makePolygon(values)

	if polygon.Contains(point){
		return true
	}else{
		return false;
	}
}


func makePolygon(values []map[string]float64) *Polygon{

	points := make([]*Point, 0);
	var point *Point
	for _, v := range values{
		point = NewPoint(v["lat"], v["lng"])
		points = append(points, point);
	}

	polygon := NewPolygon(points)

	return polygon
}


func (p *Point) UnmarshalJSON(data string) error {
	// TODO throw an error if there is an issue parsing the body.
	dec := json.NewDecoder(bytes.NewReader([]byte(data)))
	var values map[string]float64
	err := dec.Decode(&values)

	if err != nil {
		log.Print(err)
		return err
	}
	*p = *NewPoint(values["lat"], values["lng"])

	return nil
}


func (p *Point) UnmarshalJSONArray(data string) []map[string]float64 {
	// TODO throw an error if there is an issue parsing the body.

	values :=  make([]map[string]float64, 0)
	err := json.Unmarshal([]byte(data), &values)

	if err != nil {
		log.Print(err)
	}

	return values
}