package controllers

import (
	"net/http"

	"github.com/dharlequin/restful-crud-golang-app/api/responses"
)

//Home welcomes us in our API
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome User API, version 1.0.0-5")
}
