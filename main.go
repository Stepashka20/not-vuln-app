package main

import (
	"NotVulnApp/db"
	"NotVulnApp/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db.Init()
	server.Init()
}
