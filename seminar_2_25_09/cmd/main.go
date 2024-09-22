package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type HTTPError struct {
	Err        error
	HTTPStatus int
}

var (
	users  = make(map[int]*User)
	nextID = 0
)

func main() {
	http.HandleFunc("/users", handleUsersList)
	http.HandleFunc("/users/", handleUser)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Debug(err.Error())
	}
}

func handleUsersList(w http.ResponseWriter, r *http.Request) {
	var err *HTTPError
	switch r.Method {
	case "POST":
		err = createUser(w, r)
	case "GET":
		err = getUsersList(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, err.Err.Error(), err.HTTPStatus)
		return
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	var handlerErr *HTTPError

	id, err := strconv.Atoi(r.URL.Path[len("/users/"):])
	if err != nil {
		http.Error(w, "Bad argument", http.StatusBadRequest)
		http.Error(w, "Bad argument", http.StatusBadRequest)
	}
	user, exists := users[id]

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		handlerErr = getUser(w, user)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	if handlerErr != nil {
		http.Error(w, handlerErr.Err.Error(), handlerErr.HTTPStatus)
		return
	}
}

func createUser(w http.ResponseWriter, r *http.Request) *HTTPError {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return &HTTPError{
			Err:        err,
			HTTPStatus: http.StatusBadRequest,
		}
	}
	user.ID = nextID
	nextID += 1
	users[user.ID] = &user

	err = json.NewEncoder(w).Encode(users[user.ID])
	if err != nil {
		return &HTTPError{
			Err:        err,
			HTTPStatus: http.StatusInternalServerError,
		}
	}

	return nil
}

func getUsersList(w http.ResponseWriter, _ *http.Request) *HTTPError {
	userList := make([]User, 0)
	for _, user := range users {
		userList = append(userList, *user)
	}

	if err := json.NewEncoder(w).Encode(userList); err != nil {
		return &HTTPError{
			Err:        err,
			HTTPStatus: http.StatusInternalServerError,
		}
	}
	return nil
}

func getUser(w http.ResponseWriter, user *User) *HTTPError {
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return &HTTPError{
			Err:        err,
			HTTPStatus: http.StatusInternalServerError,
		}
	}
	return nil
}
