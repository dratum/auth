package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/dratum/auth/internal/repository"
	"github.com/dratum/auth/internal/repository/user"
	"github.com/dratum/auth/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051
const dbDSN = "user=postgres"

type server struct {
	auth_v1.UnimplementedAuthV1Server
	userRepository repository.UserRepository
}

func (s *server) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	user, err := s.userRepository.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &auth_v1.GetResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *server) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	id, err := s.userRepository.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &auth_v1.CreateResponse{
		Id: id,
	}, nil
}

func main() {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatal("Failed to connect database: %w", err)
	}
	defer pool.Close()

	usrRepo := user.NewRepository(pool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	auth_v1.RegisterAuthV1Server(s, &server{userRepository: usrRepo})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
