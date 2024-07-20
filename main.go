package main

import (
	"basic-trade-app/config"
	"basic-trade-app/internal/database"
	router "basic-trade-app/internal/routers"
)

var (
	PORT = "7070"
)

func main() {
	config.LoadEnv()
	database.StartDB()
	r := router.StartApp()
	r.Run(PORT)
}

