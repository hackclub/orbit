package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/fsnotify.v1"
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
`)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}
	log.SetFlags(0)

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
}

func daemonCmd(args []string) {
	fs := flag.NewFlagSet("daemon", flag.ExitOnError)
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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case e := <-watcher.Events:
				if e.Op&fsnotify.Write == fsnotify.Write {
					if err := commitEverything(); err != nil {
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

func commitEverything() error {
	if err := exec.Command("git", "add", "-A", ":/").Run(); err != nil {
		return err
	}

	if err := exec.Command("git", "commit", "-m", "", "--allow-empty-message", "--allow-empty").Run(); err != nil {
		return err
	}

	return nil
}
