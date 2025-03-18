package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dockrelix/dockrelix-agent/docker"
)

const updateURL = "http://10.1.0.73:3001/agent/update/node"

func UpdateNodeInfo(nodeInfo *docker.DockerSwarmNodeInfo) error {
	fmt.Println(nodeInfo)
	jsonData, err := json.Marshal(nodeInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal node info: %w", err)
	}

	resp, err := http.Post(updateURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
