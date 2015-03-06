// Adapted from https://code.csdn.net/flycutter/git-http-backend
package git

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/hackedu/orbit"
)

type Service struct {
	Method  string
	Handler func(handlerReq)
	RPC     string
}

type Config struct {
	ProjectRoot string
	GitBinPath  string
	UploadPack  bool
	ReceivePack bool
}

type handlerReq struct {
	w    http.ResponseWriter
	r    *http.Request
	RPC  string
	Dir  string
	File string
}

var config Config

func SetConfig(c Config) {
	config = c
}

func InitializeProject(project *orbit.Project) (projectPath string, err error) {
	projectPath = strconv.Itoa(project.ID)
	path := fmt.Sprintf("%s/%d", config.ProjectRoot, project.ID)

	if err := os.MkdirAll(path, 0744); err != nil {
		return "", err
	}

	gitCommand(path, "init", "--bare")

	return projectPath, err
}
