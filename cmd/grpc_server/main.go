package main

import (
	"context"
	"flag"
	"log"
	"net"

	authAPI "github.com/dratum/auth/internal/api/auth"
	"github.com/dratum/auth/internal/config"
	"github.com/dratum/auth/internal/repository/user"
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
	auth_v1.RegisterAuthV1Server(s, authAPI.New(authServ))

	log.Printf("server listening at %v", lis.Addr())

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
