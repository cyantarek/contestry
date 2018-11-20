package db

import (
	"log"
	"gopkg.in/mgo.v2"
)

//Encapsulation
type MongoStore struct {
	db *mgo.Database
}

//Factory Pattern
func newMgoDb(dbName string) *mgo.Database {
	conn, err := mgo.Dial("mongodb")
	if err != nil {
		log.Println(err.Error())
	}

	mgodb = conn.DB(dbName)
	return mgodb
}

//Singleton Pattern
var mgodb *mgo.Database

func GetMongoDb(dbName string) Store {
	if mgodb != nil {
		return &MongoStore{db:mgodb}
	} else {
		mgodb = newMgoDb(dbName)
		store := new(MongoStore)
		store.db = mgodb
		return store
	}
}

//Implementing or Conforming to the Interface
func (d MongoStore) ExecSQL(query string, values ...interface{}) {

}

func (d MongoStore) GetCount(query string, values ...interface{}) int {
	return 0
}

func (d *MongoStore) Migrate(models...interface{}) {

}

func (d *MongoStore) InsertData(model interface{}) {
	d.db.C("")
}

func (d *MongoStore) UpdateData(model interface{}) {

}

func (d MongoStore) GetRaw(model interface{}, query string, values...interface{}) {

}

func (d *MongoStore) GetSingle(model interface{}, query string, values interface{}) {

}

func (d MongoStore) GetAll(model interface{}) {

}
