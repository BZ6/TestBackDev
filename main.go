package main

import (
	"auth_service/config"
	"auth_service/db"
	"auth_service/router"
	"log"
	"net/http"
)

func main() {
	config.InitConfig()

	db.InitDB()

	r := router.InitRouter()

	log.Fatal(http.ListenAndServe(":8080", r))
}
