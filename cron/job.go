package cron

import (
	"context"
	"log"

	"github.com/dockrelix/dockrelix-agent/docker"
	"github.com/dockrelix/dockrelix-agent/http"

	"github.com/docker/docker/client"
	"github.com/robfig/cron/v3"
)

func StartCronJob(ctx context.Context, cli *client.Client) (*cron.Cron, error) {
	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", func() {
		nodeInfo, err := docker.GetDockerSwarmNodeInfo(ctx, cli)
		if err != nil {
			log.Printf("Failed to get Docker Swarm node info: %v", err)
			return
		}

		if err := http.UpdateNodeInfo(nodeInfo); err != nil {
			log.Printf("Failed to update node info: %v", err)
		} else {
			log.Println("Node info updated successfully")
		}
	})
	if err != nil {
		return nil, err
	}

	c.Start()
	return c, nil
}
