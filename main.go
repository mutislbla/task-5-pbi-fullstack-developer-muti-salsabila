package main

import (
	"task5-pbi/database"
	"task5-pbi/helpers"
	"task5-pbi/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	helpers.LoadEnv()
	r := gin.Default()
	database.Config()
	routes.Router(r)
}
