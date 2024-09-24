package main

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"mts2024golang/seminar_2_25_09/api"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Server struct {
	Users  map[int]*User
	NextID int

	apiv1pb.UnimplementedSeminarServiceServer
}

func NewServer() *Server {
	return &Server{
		Users:  make(map[int]*User),
		NextID: 0,
	}
}

// Получить массив с информацией о пользователях системы
// (GET /users)
func (s *Server) GetUsersList(_ context.Context, _ *apiv1pb.GetUserListRequest) (*apiv1pb.GetUserListsResponse, error) {

	userList := &apiv1pb.GetUserListsResponse{
		Users: make([]*apiv1pb.User, 0),
	}
	for _, user := range s.Users {
		addUser := &apiv1pb.User{
			Id:   int32(user.ID),
			Name: user.Name,
			Age:  strconv.Itoa(user.Age),
		}
		userList.Users = append(userList.Users, addUser)
	}

	return userList, nil
}

// Создать нового пользователя
// (POST /users)
func (s *Server) CreateUser(_ context.Context, in *apiv1pb.CreateUsersRequest) (*apiv1pb.CreateUserResponse, error) {

	user := User{
		ID:   s.NextID,
		Name: in.GetName(),
		Age:  int(in.GetAge()),
	}
	s.NextID += 1
	s.Users[user.ID] = &user

	return &apiv1pb.CreateUserResponse{
		User: &apiv1pb.User{
			Id:   int32(user.ID),
			Name: user.Name,
			Age:  strconv.Itoa(user.Age),
		},
	}, nil
}

// Получить информацию пользователя по заданному ID
// (GET /users/{id})
func (s *Server) GetUsersById(_ context.Context, in *apiv1pb.GetUsersByIdRequest) (*apiv1pb.GetUsersByIdResponse, error) {
	user, ok := s.Users[int(in.UserId)]
	if !ok {
		return nil, errors.New("User not found")
	}
	return &apiv1pb.GetUsersByIdResponse{
		User: &apiv1pb.User{
			Id:   int32(user.ID),
			Name: user.Name,
			Age:  strconv.Itoa(user.Age),
		},
	}, nil
}

func main() {

	wg := sync.WaitGroup{}
	grpcAddress := "localhost:9000"
	httpAddress := "localhost:8080"
	server := NewServer()

	grpcListen, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		slog.Debug(err.Error())
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	apiv1pb.RegisterSeminarServiceServer(grpcServer, server)

	wg.Add(1)
	go func() {
		if err := grpcServer.Serve(grpcListen); err != nil {
			slog.Debug(err.Error())
		}
		wg.Done()
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err = apiv1pb.RegisterSeminarServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		slog.Debug(err.Error())
		os.Exit(1)
	}
	httpServer := &http.Server{
		Addr:    httpAddress,
		Handler: mux,
	}

	wg.Add(1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			slog.Debug(err.Error())
		}
		wg.Done()
	}()

	wg.Wait()
}
