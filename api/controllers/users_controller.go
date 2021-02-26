package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dharlequin/restful-crud-golang-app/api/models"
	"github.com/dharlequin/restful-crud-golang-app/api/responses"
	"github.com/dharlequin/restful-crud-golang-app/api/utils/formaterror"
	"github.com/gorilla/mux"
)

//GetUsers gets all users from DB
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

//GetUser get user from DB by ID in path
func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = getUserIDFromHeaders(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}
	userFetched, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userFetched)
}

//getCurrentUserInfo returns current user info
func (server *Server) getCurrentUserInfo(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromHeaders(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}
	userFetched, err := user.FindUserByID(server.DB, uint32(id))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userFetched)
}

//UpdateUser updates existing in DB User by ID in path
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	id, err := getUserIDFromHeaders(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = checkOwnership(uid, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user.Prepare()
	err = user.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedUser)
}

//DeleteUser deletes user from DB by ID in path
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	id, err := getUserIDFromHeaders(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = checkOwnership(uid, id)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}

func getUserIDFromHeaders(r *http.Request) (uint64, error) {
	id := r.Header.Get("X-UserId")

	if id == "" {
		return 0, errors.New("User is not authorised")
	}

	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, err
	}

	return uid, nil
}

func checkOwnership(pathUserID uint64, headerUserID uint64) error {
	if pathUserID != headerUserID {
		return errors.New("You are not authorised to do that")
	}

	return nil
}