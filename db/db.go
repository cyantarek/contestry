package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/prometheus/common/log"
)

func InitDb() *gorm.DB {
	//Db, err := sql.Open("sqlite3", "db.sqlite3")
	//if err != nil {
	//	log.Println(err.Error())
	//}
	//dbmap := &gorp.DbMap{Db:Db, Dialect:gorp.SqliteDialect{}}
	//return dbmap

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}
