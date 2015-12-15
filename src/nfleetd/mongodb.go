package main

import (
	"gopkg.in/mgo.v2"
	"time"
)

const (
	MongoDBHosts = "52.77.223.167:27017"
	AuthDatabase = "nfleet"
//	AuthUserName = "guest"
//	AuthPassword = "welcome"
//	TestDatabase = "goinggo"
)


func InitMongoDB() *mgo.Session{
	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		//		Username: AuthUserName,
		//		Password: AuthPassword,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Panic("CreateSession: %s\n", err)
	}

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	mongoSession.SetMode(mgo.Monotonic, true)

	return mongoSession
}