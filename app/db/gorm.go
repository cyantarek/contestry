package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"github.com/alexsasharegan/dotenv"
	"fmt"
)

//Encapsulation
type gormDbStore struct {
	db *gorm.DB
}

//Factory Pattern
func newDb(dbName, dbType string) *gorm.DB {
	var dsn string
	if dbType == "mysql" {
		dotenv.Load()
		username := os.Getenv("MYSQL_USERNAME")
		password := os.Getenv("MYSQL_PASSWORD")
		host := os.Getenv("MYSQL_HOST")
		port := os.Getenv("MYSQL_PORT")

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)
		dbName = dsn
	}
	db, err := gorm.Open(dbType, dbName)
	if err != nil {
		log.Println(err.Error())
	}
	return db
}

//Singleton Pattern
var db *gorm.DB

func GetGormDb(dbName, dbType string) Store {
	if db != nil {
		return &gormDbStore{db: db}
	} else {
		db = newDb(dbName, dbType)
		store := new(gormDbStore)
		store.db = db
		return store
	}
}

//Implementing or Conforming to the Interface
func (d *gormDbStore) ExecSQL(query string, values ...interface{}) {
	d.db.Exec(query, values...)
}

func (d *gormDbStore) GetCount(query string, values ...interface{}) int {
	var count int
	d.db.Raw(query, values...).Count(&count)

	return count
}

func (d *gormDbStore) Migrate(models ...interface{}) {
	d.db.AutoMigrate(models...)
}

func (d *gormDbStore) InsertData(model interface{}) {
	d.db.Create(model)
}

func (d *gormDbStore) UpdateData(model interface{}) {
	d.db.Save(model)
}

func (d gormDbStore) GetRaw(model interface{}, query string, values ...interface{}) {
	d.db.Raw(query, values...).Scan(model)
}

func (d *gormDbStore) GetSingle(model interface{}, query string, values interface{}) {
	d.db.First(model, query, values)
}

func (d gormDbStore) GetAll(model interface{}) {
	d.db.Find(model)
}
