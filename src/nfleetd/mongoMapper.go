package main
import "rule"

func (m MongoSession) InsertMessageToMongoDB(msgList []rule.Message){
	go func(){
		if msgList != nil {
			for _, msg := range msgList{
				err := m.session.DB("nfleet").C("gpsdata").Insert(msg)
				if err != nil {
					log.Error("Cannot insert to Mongodb : ", err)
				}
			}
		}
	}()
}