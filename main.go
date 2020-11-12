package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	fsnotify "github.com/fsnotify/fsnotify"
)

const namespace = "main"

type stringListFlag []string

func (v *stringListFlag) Set(value string) error {
	*v = append(*v, value)
	return nil
}

func (v *stringListFlag) String() string {
	return fmt.Sprint(*v)
}

var (
	command    = flag.String("command", "", "the command to execute when a config is updated")
	args       stringListFlag
	volumeDirs stringListFlag
)

func main() {
	flag.Var(&volumeDirs, "volume-dir", "the config map volume directory to watch for updates; may be used multiple times")
	flag.Var(&args, "argument", "an argument for the command; may be used multiple times; optional")
	flag.Parse()

	if len(volumeDirs) < 1 {
		log.Println("Missing --volume-dir=")
		log.Println()
		flag.Usage()
		os.Exit(1)
	}

	if len(*command) < 1 {
		log.Println("Missing --command=")
		log.Println()
		flag.Usage()
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("config map updated: ", event)

				cmd := exec.Command(*command, args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()

				if err != nil {
					log.Println("command failed:", err)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	for _, d := range volumeDirs {
		log.Printf("Watching directory: %q", d)
		err = watcher.Add(d)
		if err != nil {
			log.Fatal(err)
		}
	}

	done := make(chan bool)
	<-done // Block forever
}
