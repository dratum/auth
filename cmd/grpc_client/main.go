package main

import (
	"context"
	"log"
	"time"

	desc "github.com/dratum/auth/pkg/auth_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
	userId  = 2
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewAuthV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: userId})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}

	log.Printf(color.RedString("Auth info:\n"), color.GreenString("%+v", r.GetRole()))
}
