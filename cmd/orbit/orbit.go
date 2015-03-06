package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/hackedu/orbit"
	"gopkg.in/fsnotify.v1"
)

var (
	baseURLStr = flag.String("url", "http://mew.hackedu.us:4000", "base URL of orbit")
	baseURL    *url.URL

	apiClient = orbit.NewClient(nil)
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `orbit puts your development environment in the cloud.

Usage:

    orbit [options] command [arg...]

The commands are:
`)
		for _, c := range subcmds {
			fmt.Fprintf(os.Stderr, "    %-24s %s\n", c.name, c.description)
		}
		fmt.Fprintln(os.Stderr, `
Use "orbit command -h" for more information about a command.

The options are:
`)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}
	log.SetFlags(0)

	var err error
	baseURL, err = url.Parse(*baseURLStr)
	if err != nil {
		log.Fatal(err)
	}
	apiClient.BaseURL = baseURL.ResolveReference(&url.URL{Path: "/api/"})

	subcmd := flag.Arg(0)
	for _, c := range subcmds {
		if c.name == subcmd {
			c.run(flag.Args()[1:])
			return
		}
	}

	fmt.Fprintf(os.Stderr, "unknown subcmd %q\n", subcmd)
	fmt.Fprintln(os.Stderr, `Run "orbit -h" for usage.`)
	os.Exit(1)
}

type subcmd struct {
	name        string
	description string
	run         func(args []string)
}

var subcmds = []subcmd{
	{"daemon", "start the orbit daemon", daemonCmd},
	{"create-project", "create a new project", createProjectCmd},
	{"create-service", "create a new service", createServiceCmd},
	{"run", "run command in current prject", runCommandCmd},
	{"open", "open a project in your browser", openCmd},
}

func daemonCmd(args []string) {
	fs := flag.NewFlagSet("daemon", flag.ExitOnError)
	rateLimit := fs.Float64("q", 0, "rate limit in seconds (0 to disable)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: orbit daemon [options]

Start the Orbit daemon that watches for and acts on file changes.

The options are:
`)
		fs.PrintDefaults()
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 0 {
		fs.Usage()
	}

	// TODO: Fix throttle
	var throttle <-chan time.Time
	if *rateLimit > 0 {
		throttle = time.Tick(time.Duration(1e6/(*rateLimit)) * time.Microsecond)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			if *rateLimit > 0 {
				<-throttle
			}

			select {
			case e := <-watcher.Events:
				if e.Op&fsnotify.Write == fsnotify.Write {
					if err := commitAndPushEverything(); err != nil {
						log.Fatal("error committing changes")
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func commitAndPushEverything() error {
	if err := exec.Command("git", "add", "-A", ":/").Run(); err != nil {
		return err
	}

	if err := exec.Command("git", "commit", "-m", "", "--allow-empty-message", "--allow-empty").Run(); err != nil {
		return err
	}

	if err := exec.Command("git", "push").Run(); err != nil {
		return err
	}

	return nil
}

func createProjectCmd(args []string) {
	fs := flag.NewFlagSet("create-project", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: orbit create-project [project name] [options]

Create a new project on Orbit.
`)
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 1 {
		fs.Usage()
	}

	projectName := fs.Args()[0]

	var project orbit.Project
	if err := apiClient.Projects.Create(&project); err != nil {
		log.Fatal(err)
	}
	cloneURL := baseURL.ResolveReference(&url.URL{Path: "/git/"}).ResolveReference(&url.URL{Path: project.GitPath})
	if err := exec.Command("git", "clone", cloneURL.String(), projectName).Run(); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(
		fmt.Sprintf("%s/.orbitrc", projectName),
		[]byte(strconv.Itoa(project.ID)),
		0644,
	); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s created successfully\n", projectName)
}

func createServiceCmd(args []string) {
	fs := flag.NewFlagSet("create-service", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: orbit create-service [service type] [options]

Create a new service for the current project on Orbit.
`)
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 1 {
		fs.Usage()
	}

	serviceType := fs.Args()[0]

	orbitrc, err := ioutil.ReadFile(".orbitrc")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "You must be in an orbit project to create a service.")
			os.Exit(1)
		}

		log.Fatal(err)
	}

	projectID, err := strconv.Atoi(string(orbitrc))
	if err != nil {
		fmt.Fprintln(os.Stderr, ".orbitrc is corrupted. Please fix it and try again.")
		os.Exit(1)
	}

	service := func() *orbit.Service {
		switch serviceType {
		case "rails":
			return &orbit.Service{
				Type:        serviceType,
				ProjectID:   projectID,
				PortExposed: "3000",
			}
		case "node":
			return &orbit.Service{
				Type:        serviceType,
				ProjectID:   projectID,
				PortExposed: "3000",
			}
		case "tpires/neo4j":
			return &orbit.Service{
				Type:        serviceType,
				ProjectID:   projectID,
				PortExposed: "7474",
			}
		}
		return nil
	}()
	if service == nil {
		fmt.Fprintf(os.Stderr, "%s is not a valid project type.\n", serviceType)
		os.Exit(1)
	}

	if err := apiClient.Services.Create(service); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s service created successfully\n", serviceType)
}

func runCommandCmd(args []string) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: orbit run [container type] [command] [options]

Run a command in the current project on Orbit.
`)
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() < 2 {
		fs.Usage()
	}

	orbitrc, err := ioutil.ReadFile(".orbitrc")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "You must be in an orbit project to create a service.")
			os.Exit(1)
		}

		log.Fatal(err)
	}

	projectID, err := strconv.Atoi(strings.TrimSpace(string(orbitrc)))
	if err != nil {
		fmt.Fprintln(os.Stderr, ".orbitrc is corrupted. Please fix it and try again.")
		os.Exit(1)
	}

	cmd := orbit.ProjectCmd{
		ContainerType: fs.Args()[0],
		Command:       fs.Args()[1:],
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(cmd); err != nil {
		log.Fatal(err)
	}

	u := apiClient.BaseURL.ResolveReference(&url.URL{
		Path: fmt.Sprintf("/api/projects/%d/run_command", projectID),
	})
	resp, err := http.Post(u.String(), "application/json", buf)
	if err != nil {
		log.Fatal("error running command: could not make http request")
	}

	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)

	if output, err := exec.Command("git", "pull").Output(); err != nil {
		fmt.Println(string(output))
		log.Fatal("error pulling file changes from command:", err)
	}
}

func openCmd(args []string) {
	fs := flag.NewFlagSet("open", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: orbit open [service] [options]

Opens the provided service in the browser.
`)
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 1 {
		fs.Usage()
	}

	serviceType := fs.Args()[0]

	orbitrc, err := ioutil.ReadFile(".orbitrc")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "You must be in an Orbit project to run this.")
			os.Exit(1)
		}

		log.Fatal(err)
	}

	projectID, err := strconv.Atoi(string(orbitrc))
	if err != nil {
		fmt.Fprintln(os.Stderr, ".orbitrc is corrupted. Please fix it and try again.")
		os.Exit(1)
	}

	services, err := apiClient.Services.List(projectID)
	if err != nil {
		log.Fatal(err)
	}

	service := func() *orbit.Service {
		for _, service := range services {
			if service.Type == serviceType {
				return service
			}
		}
		return nil
	}()
	if service == nil {
		fmt.Fprintf(os.Stderr, "the %s service does not exist for this project", serviceType)
		os.Exit(1)
	}

	cmd := exec.Command("open", fmt.Sprintf("http://mewtwo.hackedu.us:%s", service.HostPort))
	if err := cmd.Run(); err != nil {
		log.Fatal("error opening project in browser:", err)
	}
}
