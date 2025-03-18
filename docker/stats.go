package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/dockrelix/dockrelix-agent/utils"

	"github.com/docker/docker/client"
)

func GetDockerSwarmNodeInfo(ctx context.Context, cli *client.Client) (*DockerSwarmNodeInfo, error) {
	info, err := cli.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker info: %w", err)
	}

	version, err := cli.ServerVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker version: %w", err)
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var totalUsedMemory, totalUsedCPU float64
    for _, container := range containers {
        stats, err := cli.ContainerStats(ctx, container.ID, false)
        if err != nil {
            log.Printf("Failed to get stats for container %s: %v", container.ID, err)
            continue
        }
        defer stats.Body.Close()

        var statsJSON map[string]interface{}
        if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
            log.Printf("Failed to decode stats for container %s: %v", container.ID, err)
            continue
        }

        memoryStats, ok := statsJSON["memory_stats"].(map[string]interface{})
        if !ok {
            log.Printf("Invalid memory stats for container %s", container.ID)
            continue
        }

        memoryUsage, ok := memoryStats["usage"].(float64)
        if !ok {
            log.Printf("Invalid memory usage for container %s", container.ID)
            continue
        }

        cpuStats, ok := statsJSON["cpu_stats"].(map[string]interface{})
        if !ok {
            log.Printf("Invalid CPU stats for container %s", container.ID)
            continue
        }

        cpuUsage, ok := cpuStats["cpu_usage"].(map[string]interface{})
        if !ok {
            log.Printf("Invalid CPU usage for container %s", container.ID)
            continue
        }

        totalUsage, ok := cpuUsage["total_usage"].(float64)
        if !ok {
            log.Printf("Invalid total CPU usage for container %s", container.ID)
            continue
        }

        memoryUsageGiB := utils.BytesToGiB(uint64(memoryUsage))
        totalUsedMemory += memoryUsageGiB
        totalUsedCPU += totalUsage
    }

	nodeInfo := &DockerSwarmNodeInfo{
		NodeID:           info.Swarm.NodeID,
		Hostname:         info.Name,
		IPAddress:        info.Swarm.NodeAddr,
		Role:             "worker",
		Platform:         info.OSType,
		Arch:             info.Architecture,
		CPU:              info.NCPU,
		Memory:           utils.BytesToGiB(uint64(info.MemTotal)),
		UsedMemory:       totalUsedMemory,
		UsedCPU:          totalUsedCPU,
		KernelVersion:    info.KernelVersion,
		OperatingSystem:  info.OperatingSystem,
		DockerVersion:    version.Version,
	}

	if info.Swarm.ControlAvailable {
		nodeInfo.Role = "manager"
	}

	return nodeInfo, nil
}
