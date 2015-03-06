package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/hackedu/orbit/api"
	"github.com/hackedu/orbit/datastore"
	"github.com/hackedu/orbit/docker"
	"github.com/hackedu/orbit/git"
)

var (
	baseURLStr = flag.String("url", "http://mew.hackedu.us:4000", "base URL of orbit")
	baseURL    *url.URL
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `orbit-server is orbit's server compontent.

Usage:

    orbit [options] command [arg...]

The commands are:
`)
		for _, c := range subcmds {
			fmt.Fprintf(os.Stderr, "    %-24s %s\n", c.name, c.description)
		}
		fmt.Fprintln(os.Stderr, `
Use "orbit-server command -h" for more information about a command.

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
	docker.BaseURL = baseURL

	subcmd := flag.Arg(0)
	for _, c := range subcmds {
		if c.name == subcmd {
			c.run(flag.Args()[1:])
			return
		}
	}

	fmt.Fprintf(os.Stderr, "unknown subcmd %q\n", subcmd)
	fmt.Fprintln(os.Stderr, `Run "orbit-server -h" for usage.`)
	os.Exit(1)
}

type subcmd struct {
	name        string
	description string
	run         func(args []string)
}

var subcmds = []subcmd{
	{"serve", "start web server", serveCmd},
	{"createdb", "create the database schema", createDBCmd},
}

func serveCmd(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	httpAddr := fs.String("http", ":5000", "HTTP service address")
	gitProjectsRoot := fs.String("git-root", ".", "git projects root")
	gitBin := fs.String("git-bin", "/usr/bin/git", "path to git binary")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: orbit-server serve [options]

Starts the web server that serves the API.

The options are:
`)
		fs.PrintDefaults()
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 0 {
		fs.Usage()
	}

	git.SetConfig(git.Config{
		ProjectRoot: *gitProjectsRoot,
		GitBinPath:  *gitBin,
		UploadPack:  true,
		ReceivePack: true,
	})

	datastore.Connect()

	m := http.NewServeMux()
	m.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	m.Handle("/git/", http.StripPrefix("/git", git.Handler()))

	log.Print("Listening on ", *httpAddr)
	err := http.ListenAndServe(*httpAddr, m)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func createDBCmd(args []string) {
	fs := flag.NewFlagSet("createdb", flag.ExitOnError)
	drop := fs.Bool("drop", false, "drop DB before creating")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: orbit-server createdb [options]

Creates the necessary DB tables and indexes.

The options are:
`)
		fs.PrintDefaults()
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 0 {
		fs.Usage()
	}

	datastore.Connect()
	if *drop {
		datastore.Drop()
	}
	datastore.Create()
}
