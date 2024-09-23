package main

import (
	"encoding/json"
	"log/slog"
	"mts2024golang/seminar_2_25_09/api"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Server struct {
	Users  map[int]*User
	NextID int
}

func NewServer() *Server {
	return &Server{
		Users:  make(map[int]*User),
		NextID: 0,
	}
}

// Получить массив с информацией о пользователях системы
// (GET /users)
func (s *Server) GetUsersList(w http.ResponseWriter, _ *http.Request) {
	userList := make([]User, 0)
	for _, user := range s.Users {
		userList = append(userList, *user)
	}

	if err := json.NewEncoder(w).Encode(userList); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Создать нового пользователя
// (POST /users)
func (s *Server) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad argument", http.StatusBadRequest)
	}
	user.ID = s.NextID
	s.NextID += 1
	s.Users[user.ID] = &user

	err = json.NewEncoder(w).Encode(s.Users[user.ID])
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Получить информацию пользователя по заданному ID
// (GET /users/{id})
func (s *Server) GetUserByID(w http.ResponseWriter, _ *http.Request, id int) {
	user, ok := s.Users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func main() {
	server := NewServer()

	r := http.NewServeMux()

	h := api.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	if err := s.ListenAndServe(); err != nil {
		slog.Debug(err.Error())
	}
}
