package docker

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hackedu/orbit"
)

var (
	Store   = orbit.NewClient(nil)
	BaseURL *url.URL
)

func CommandInContainer(containerID string, command ...string) *exec.Cmd {
	return exec.Command("docker", "exec",
		containerID,
		"/bin/sh", "-c", "cd /usr/src/app && "+strings.Join(command, " "),
	)
}

func UpdateProjectFilesInServicesForProject(projectID int) error {
	services, err := Store.Services.List(projectID)
	if err != nil {
		return err
	}

	for _, service := range services {
		cmd := exec.Command("docker", "exec",
			service.ContainerID,
			"/bin/sh", "-c", "cd /usr/src/app && git pull",
		)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func RunContainer(projectID int, image, port string) (containerID, hostPort string, err error) {
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

	if err := cloneProjectRepo(containerID, projectID); err != nil {
		return "", "", err
	}

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

func cloneProjectRepo(containerID string, projectID int) error {
	gitURL := BaseURL.ResolveReference(&url.URL{Path: "/git/"})
	gitURL = gitURL.ResolveReference(&url.URL{Path: strconv.Itoa(projectID)})
	fmt.Println(gitURL)
	cmd := exec.Command("docker", "exec",
		containerID,
		"git", "clone", gitURL.String(), "/usr/src/app",
	)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
