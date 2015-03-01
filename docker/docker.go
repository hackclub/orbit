package docker

import (
	"encoding/json"
	"os/exec"
	"strings"
)

func RunContainer(image, port string) (containerID, hostPort string, err error) {
	output, err := exec.Command("docker", "run",
		"-d",
		"-p", port,
		image,
		"/bin/sh", "-c", "while true; do sleep 1; done",
	).Output()
	if err != nil {
		return "", "", err
	}
	containerID = strings.TrimSpace(string(output))

	output, err = exec.Command("docker", "inspect", containerID).Output()
	if err != nil {
		return "", "", err
	}

	var inspectOutput []struct {
		NetworkSettings struct {
			Ports map[string][]struct {
				HostPort string
			}
		}
	}
	if err := json.Unmarshal(output, &inspectOutput); err != nil {
		return "", "", err
	}

	hostPort = inspectOutput[0].NetworkSettings.Ports[port+"/tcp"][0].HostPort

	return containerID, hostPort, nil
}
