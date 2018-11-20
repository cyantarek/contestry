package main

import (
	"net/http"
	"time"
	"encoding/gob"

	"./app/routes"
	"./app/models"
	"github.com/alexsasharegan/dotenv"
	"os"
)

/*
Two features

01. Dependency injection based controllers. Composing Handler structs wrapping dependencies. Create Controller
Functions as Methods on that struct
02. Interface driven package design and 3rd party dependencies
*/

func main() {
	dotenv.Load()
	addr := ":" + os.Getenv("PORT")
	server := &http.Server{
		Addr:addr,
		Handler: routes.GetRouter(),
		WriteTimeout:time.Second*5,
		ReadTimeout:time.Second*5,
	}

	gob.Register(models.User{})
	//file, _ := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	//log.SetOutput(file)
	//log.SetFlags(0)
	server.ListenAndServe()
}
