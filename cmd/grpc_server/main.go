package main

import (
	"context"
	"flag"
	"log"

	"github.com/dratum/auth/internal/app"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
