package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/dratum/auth/pkg/auth_v1"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051
const dbDSN = ""

type server struct {
	auth_v1.UnimplementedAuthV1Server
	db *pgxpool.Pool
}

func ConnectDB(ctx context.Context, dbDSN string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	ctx := context.Background()

	con, err := ConnectDB(ctx, dbDSN)
	if err != nil {
		log.Fatal("Failed to connect database: %w", err)
	}
	defer con.Close(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	auth_v1.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
