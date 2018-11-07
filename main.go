package main

import (
	"./server"
	"./router"
	"./db"
	"./user"
	"./contest"
)

func main() {
	//Migrate()
	//server := http.Server{
	//	Handler:GetRoutes(),
	//	Addr:":50051",
	//}
	//log.Println("Server started...")
	//err := server.ListenAndServe()
	//if err != nil {
	//	log.Println(err.Error())
	//}
	Db := db.InitDb()
	Db.AutoMigrate(user.User{}, contest.Contest{})
	defer Db.Close()
	server.Run("50052", router.GetRouter())
}
