package git

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hackedu/orbit/docker"
)

func serviceRPC(hr handlerReq) {
	w, r, rpc, dir := hr.w, hr.r, hr.RPC, hr.Dir
	access := hasAccess(r, dir, rpc, true)

	if access == false {
		renderNoAccess(w)
		return
	}

	input, _ := ioutil.ReadAll(r.Body)

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", rpc))
	w.WriteHeader(http.StatusOK)

	args := []string{rpc, "--stateless-rpc", dir}
	cmd := exec.Command(config.GitBinPath, args...)
	cmd.Dir = dir
	in, err := cmd.StdinPipe()
	if err != nil {
		log.Print(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Print(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Print(err)
	}

	in.Write(input)
	io.Copy(w, stdout)
	cmd.Wait()

	if hr.RPC == "receive-pack" {
		split := strings.Split(hr.Dir, "/")
		projectID, err := strconv.Atoi(split[len(split)-1])
		if err != nil {
			log.Print(err)
			return
		}

		if err := docker.UpdateProjectFilesInServicesForProject(projectID); err != nil {
			log.Print(err)
			return
		}
	}
}

func getInfoRefs(hr handlerReq) {
	w, r, dir := hr.w, hr.r, hr.Dir
	service_name := getServiceType(r)
	access := hasAccess(r, dir, service_name, false)

	if access {
		args := []string{service_name, "--stateless-rpc", "--advertise-refs", "."}
		refs := gitCommand(dir, args...)

		hdrNoCache(w)
		w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", service_name))
		w.WriteHeader(http.StatusOK)
		w.Write(packetWrite("# service=git-" + service_name + "\n"))
		w.Write(packetFlush())
		w.Write(refs)
	} else {
		updateServerInfo(dir)
		hdrNoCache(w)
		sendFile("text/plain; charset=utf-8", hr)
	}
}

func getInfoPacks(hr handlerReq) {
	hdrCacheForever(hr.w)
	sendFile("text/plain; charset=utf-8", hr)
}

func getLooseObject(hr handlerReq) {
	hdrCacheForever(hr.w)
	sendFile("application/x-git-loose-object", hr)
}

func getPackFile(hr handlerReq) {
	hdrCacheForever(hr.w)
	sendFile("application/x-git-packed-objects", hr)
}

func getIdxFile(hr handlerReq) {
	hdrCacheForever(hr.w)
	sendFile("application/x-git-packed-objects-toc", hr)
}

func getTextFile(hr handlerReq) {
	hdrNoCache(hr.w)
	sendFile("text/plain", hr)
}
