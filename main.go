package main

import (
	"workshoptdd/config"
	"workshoptdd/routes"
)

func main() {
	db := config.InitDatabase("root:root@tcp(127.0.0.1:3306)/go_tdd?charset=utf8mb4&parseTime=True&loc=Local")
	app := routes.InitRoutes(db)

	app.Run(":8000")
}
