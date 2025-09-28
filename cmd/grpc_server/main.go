package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/dratum/auth/internal/config"
	"github.com/dratum/auth/internal/converter"
	"github.com/dratum/auth/internal/repository/user"
	"github.com/dratum/auth/internal/service"
	"github.com/dratum/auth/internal/service/auth"
	"github.com/dratum/auth/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	auth_v1.UnimplementedAuthV1Server
	authService service.AuthService
}

func (s *server) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	id, err := s.authService.Create(ctx, converter.ToUserFromAuth(req))
	if err != nil {
		return nil, err
	}

	return &auth_v1.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	user, err := s.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	roleValue, exists := auth_v1.Role_value[user.Role]
	if !exists {
		return nil, fmt.Errorf("invalid role value from database: %s", user.Role)
	}
	return converter.ToUserFromService(user, roleValue), nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGrpcConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get postgres config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer pool.Close()

	usrRepo := user.NewRepository(pool)
	authServ := auth.NewService(usrRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	auth_v1.RegisterAuthV1Server(s, &server{authService: authServ})

	log.Printf("server listening at %v", lis.Addr())

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
