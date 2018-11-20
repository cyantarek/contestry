package db

import (

)

//Abstraction - Services exposed by db package
type Store interface {
	ExecSQL(query string, values ...interface{})
	GetCount(query string, values ...interface{}) int
	Migrate(models...interface{})
	InsertData(model interface{})
	UpdateData(model interface{})
	GetRaw(model interface{}, query string, values...interface{})
	GetSingle(model interface{}, query string, values interface{})
	GetAll(model interface{})
}

/*
Package Db is a 3rd party dependency

package db publishes a public interface "Store" with public methods. And a public function to create
the database the public function returns private gormDbStore, an struct that composes the database driver
so other package can not access it's details. The gormDbStore struct implements all the methods of the
Store interface.

Interface is a written specification that defines what sort of functionality any component should have.
It's like architecture specification. Now any component can implement the architecture.

So, other package does not know how the database is implemented.

They just declared a variable for Store interface, then call the public function to create database. The
returned private database will be assigned to the Store interface variable, as it implements the interface.

This way we can easily swap out one database implements with another one satisfying the Store interface, this
way our code will not break.
*/
