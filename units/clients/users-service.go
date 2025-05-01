package clients

import (
	"log"
	"sync"

	usersv1 "github.com/QuizWars-Ecosystem/load-tests/gen/external/users/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client  *UsersService
	once    sync.Once
	initErr error
)

type UsersService struct {
	usersv1.UsersAuthServiceClient
	usersv1.UsersProfileServiceClient
	usersv1.UsersSocialServiceClient
	usersv1.UsersAdminServiceClient
}

func GetUsersClient(addr string) (*UsersService, error) {
	once.Do(func() {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			initErr = err
			log.Fatalf("failed client connection: %v", err)
		}

		client = &UsersService{
			UsersAuthServiceClient:    usersv1.NewUsersAuthServiceClient(conn),
			UsersProfileServiceClient: usersv1.NewUsersProfileServiceClient(conn),
			UsersSocialServiceClient:  usersv1.NewUsersSocialServiceClient(conn),
			UsersAdminServiceClient:   usersv1.NewUsersAdminServiceClient(conn),
		}
	})

	return client, initErr
}
