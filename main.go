package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// User model
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// Users var imitetes a DB list of Users
var Users []User

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", returnAllUsers)
	myRouter.HandleFunc("/user", createNewUser).Methods("POST")
	myRouter.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{id}", returnSingleUser)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Users = []User{
		{ID: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@gmail.com", Phone: "532532"},
		{ID: 2, FirstName: "Tony", LastName: "P", Email: "tony.p@gmail.com", Phone: "43363453"},
	}
	handleRequests()
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: returnAllArticles")
	json.NewEncoder(w).Encode(Users)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromPath(w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	for _, user := range Users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
		}
	}
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprint(w, "Error reading request body")
		return
	}

	var user User
	json.Unmarshal(reqBody, &user)

	if validateUserFields(user) {
		user.ID = generateNewUserID()

		Users = append(Users, user)

		json.NewEncoder(w).Encode(user)
	} else {
		fmt.Fprintln(w, "First Name and Last Name cannot be empty!")
	}
}

func validateUserFields(user User) bool {
	if user.FirstName == "" || user.LastName == "" {
		return false
	}

	return true
}

func generateNewUserID() int {
	id := 0

	for _, user := range Users {
		if user.ID > id {
			id = user.ID
		}
	}

	return id + 1
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromPath(w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	for index, user := range Users {
		if user.ID == id {
			Users = append(Users[:index], Users[index+1:]...)
		}
	}

	json.NewEncoder(w).Encode(Users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := getUserIDFromPath(w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprint(w, "Error reading request body", err)
		return
	}

	var newUser User
	json.Unmarshal(reqBody, &newUser)

	if validateUserFields(newUser) {
		for index, user := range Users {
			if user.ID == id {
				newUser.ID = id
				Users[index] = newUser
			}
		}

		json.NewEncoder(w).Encode(Users)
	} else {
		fmt.Fprintln(w, "First Name and Last Name cannot be empty!")
	}
}

func getUserIDFromPath(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	key := vars["id"]
	id, err := strconv.Atoi(key)

	fmt.Println("Key: " + key)

	if err != nil {
		return 0, fmt.Errorf("Entered ID is not a valid number")
	}

	return id, nil
}
