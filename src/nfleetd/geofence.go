package main

import (
	"encoding/json"
	"bytes"
	"rule"
	"fmt"
)

//func checkInOutGeofence(geofenceList []Geofence, msgList []rule.Message) []Geofence{
func CheckInOutGeofence(msgList []rule.Message, geofenceList []Geofence, mysqlSession *SQLMapper){


//	var IMEI string
//	for _, msg := range msgList {
//		IMEI = msg.IMEI
//	}
	//select geofence state
//	mysqlSession.GetGeofenceState(IMEI, geofenceAndStmt)

	resultList := Checkpoint(geofenceList, msgList)

	if resultList != nil{
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
						fmt.Println("1")
					}else{
						fmt.Println("2")
					}
				}else{
					updateToOutState(geofence,  msg.IMEI)
					fmt.Println("3")
				}
			}else if geofence.shape == 2{
				if PointsWithinPolygon(geofence, point) {
					if CheckInoutState(geofence, msg.IMEI){
						resultList = append(resultList, geofence)
						geofence.state[msg.IMEI] = 2
						fmt.Println("4")
					}else{
						fmt.Println("5")
					}
				}else{
					updateToOutState(geofence,  msg.IMEI)
					fmt.Println("6")
				}
			}
		}
	}

	return resultList
}


func CheckInoutState(geofence Geofence, imei string) bool{

	flag, _ := geofence.state[imei]

	if flag == 1 || flag == 0{
		fmt.Println("7")
		 return true
	}else if flag ==2 {
		fmt.Println("8")
		return false
	}else{
		fmt.Println("9")
		return false
	}
}

func updateToOutState(geofence Geofence,  imei string){
	flag := geofence.state[imei]

	if flag == 2{
		fmt.Println("10")
		flag = 1
	}
	fmt.Println("11")
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