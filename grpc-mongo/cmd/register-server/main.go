package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	module "grpc-mongo/modules"
	userpb "grpc-mongo/protos/v1/user"
)

const portNumber = "9000"

type userService struct {
	userpb.UserServiceServer
}

func (s *userService) SignUp(ctx context.Context, req *userpb.SignUpRequest) (*userpb.SignUpResponse, error) {
	user_id, user_name, user_gender, user_email := req.User.Id, req.User.Name, req.User.Gender, req.User.Email
	user_password, _ := module.Hash(req.Password)
	newUser := module.User{user_id, user_name, user_gender, user_email, user_password}

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
		return &userpb.SignUpResponse{
			Success: false,
			Message: ".env file not gound",
		}, nil
	}
	MONGODB_URI := os.Getenv("MONGODB_URI")

	module.MongoConnection(MONGODB_URI)

	client := module.MongoConnection(MONGODB_URI)
	defer module.MongoDisconnection(client)

	module.MongoUserInsert(newUser, client.Database("grpc").Collection("users"))

	return &userpb.SignUpResponse{
		Success: true,
		Message: "User Seccessfully Inserted",
	}, nil

}

// func (s *userService) LogIn(ctx context.Context, req *userpb.LogInRequest) (*userpb.LogInResponse, error) {
// 	user_id, user_password := req.Id, req.Password
// 	logIn := module.LogIn{user_id, user_password}

// 	if err := godotenv.Load(); err != nil {
// 		log.Println(".env file not found")
// 		return &userpb.LogInResponse{
// 			Success: false,
// 			User:    &userpb.User{},
// 		}, nil
// 	}
// 	MONGODB_URI := os.Getenv("MONGODB_URI")

// 	module.MongoConnection(MONGODB_URI)

// 	client := module.MongoConnection(MONGODB_URI)
// 	defer module.MongoDisconnection(client)

// 	module.MongoUserLogInFind(logIn, client.Database("grpc").Collection("users"))

// 	return &userpb.LogInResponse{
// 		Success: true,
// 		Message: "",
// 		User: &userpb.User{
// 			Id:     "1",
// 			Name:   "John",
// 			Gender: "Male",
// 			Email:  "john@example.com",
// 		},
// 	}, nil

// }

func main() {
	lis, err := net.Listen("tcp", ":"+portNumber)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &userService{})

	log.Printf("start gRPC server on %s port", portNumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}