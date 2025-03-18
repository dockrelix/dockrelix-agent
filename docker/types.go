package docker

type ContainerInfo struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	State      string  `json:"state"`
	Status     string  `json:"status"`
	Image      string  `json:"image"`
	UsedMemory float64 `json:"usedMemory"`
	UsedCPU    float64 `json:"usedCpu"`
}

type DockerSwarmNodeInfo struct {
	NodeID           string  `json:"node_id"`
	Hostname         string  `json:"hostname"`
	IPAddress        string  `json:"ip_address"`
	Role             string  `json:"role"`
	Platform         string  `json:"platform"`
	Arch             string  `json:"arch"`
	CPU              int     `json:"cpu"`
	Memory           float64 `json:"memory"`
	UsedMemory       float64 `json:"usedMemory"`
	UsedCPU          float64 `json:"usedCpu"`
	KernelVersion    string  `json:"kernelVersion"`
	OperatingSystem  string  `json:"operatingSystem"`
	DockerVersion    string  `json:"dockerVersion"`
}
