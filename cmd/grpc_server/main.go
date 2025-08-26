package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/dratum/auth/internal/domain/models"
	"github.com/dratum/auth/pkg/auth_v1"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051
const dbDSN = ""

type server struct {
	auth_v1.UnimplementedAuthV1Server
	db *pgx.Conn
}

func ConnectDB(ctx context.Context, dbDSN string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *server) Get(
	ctx context.Context,
	req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	query := `
			  select id
			       , name
			       , email
			       , role
			       , created_at
			       , updated_at
			  from users u
			  where u.id = $1
			`

	var user models.User
	var roleStr string

	err := s.db.QueryRow(ctx, query, req.Id).Scan(
		&user.Id,
		&user.Email,
		&user.Name,
		&roleStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	roleValue, exists := auth_v1.Role_value[strings.ToUpper(roleStr)]
	if !exists {
		return nil, fmt.Errorf("invalid role value from database: %s", roleStr)
	}

	return &auth_v1.GetResponse{
		Id:        user.Id,
		Email:     user.Email,
		Name:      user.Name,
		Role:      auth_v1.Role(roleValue),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (s *server) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	return &auth_v1.CreateResponse{}, nil
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
