package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	_ "github.com/swaggo/swag/example/celler/httputil"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// @title API Примера для семинара
// @version 0.1
// @description Пhимер описания API демонстрации на семинаре

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/", handleUser)
	http.Handle("/swagger/*", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Debug(err.Error())
	}
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
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

// createUser godoc
// @Summary Создать нового пользователя
// @Description Создать нового пользователя
// @Accept  json
// @Produce  json
// @Param   user	body	User	true	"Данные пользователя"
// @Success 200 {object} User
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /users [post]
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

// getUserList godoc
// @Summary Получить список пользователей
// @Description Получить массив с информацией о пользователях системы
// @Accept  json
// @Produce  json
// @Success 200 {array} User
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /users [get]
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

// getUser godoc
// @Summary Вывести информацию поп пользователю
// @Description Получить информацию пользователя по заданному ID
// @Accept  json
// @Produce  json
// @Param   id	path		int		true	"ID пользователя"
// @Success 200 {object} User
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /users/{id} [get]
func getUser(w http.ResponseWriter, user *User) *HTTPError {
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return &HTTPError{
			Err:        err,
			HTTPStatus: http.StatusInternalServerError,
		}
	}
	return nil
}
