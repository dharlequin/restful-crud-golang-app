package api

import (
	"os"

	"github.com/dharlequin/restful-crud-golang-app/api/controllers"
	"github.com/dharlequin/restful-crud-golang-app/api/seed"
)

var server = controllers.Server{}

//Run runs all server parts
func Run() {
	server.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	seed.Load(server.DB)

	server.Run(":8080")
}
