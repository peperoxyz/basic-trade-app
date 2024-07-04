package main

import (
	"basic-trade-app/config"
	"basic-trade-app/internal/database"
	router "basic-trade-app/internal/routers"
	"os"
)

var (
	PORT = ":" + os.Getenv("PORT")
)

func main() {
	config.LoadEnv()
	database.StartDB()
	r := router.StartApp()
	r.Run(PORT)
}

