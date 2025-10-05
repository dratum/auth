package app

import (
	"context"
	"log"

	"github.com/dratum/auth/internal/api/auth"
	"github.com/dratum/auth/internal/closer"
	"github.com/dratum/auth/internal/config"
	"github.com/dratum/auth/internal/repository"
	"github.com/dratum/auth/internal/repository/user"
	"github.com/dratum/auth/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgPool         *pgxpool.Pool
	userRepository repository.UserRepository

	authService service.AuthService

	server *auth.Server
}

func New() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {

		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get postgres config: %v", err)

		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGrpcConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)

		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) PGPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.Connect(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("Failed to connect database: %v", err)

		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("Ping error: %v", err)
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.pgPool = pool
	}

	return s.pgPool
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = user.NewRepository(s.PGPool(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = s.UserRepository(ctx)
	}

	return s.authService
}

func (s *serviceProvider) AuthServer(ctx context.Context) *auth.Server {
	if s.server == nil {
		s.server = auth.New(s.AuthService(ctx))
	}

	return s.server
}
