package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dockrelix/dockrelix-agent/cron"
	"github.com/dockrelix/dockrelix-agent/docker"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := docker.NewDockerClient()
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	cronJob, err := cron.StartCronJob(ctx, cli)
	if err != nil {
		log.Fatalf("Failed to start cron job: %v", err)
	}

	log.Println("started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down gracefully...")
	cronJob.Stop()
	cancel()
}
