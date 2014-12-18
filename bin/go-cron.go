package main

import (
	//"fmt"
	gocron "github.com/odise/go-cron"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var build string

func main() {
	log.Println("Running version: %s", build)

	if len(os.Args) < 3 {
		log.Fatalf("run: go-cron <schedule> <command>")
	}

	c, wg := gocron.Create()

	go gocron.Start(c)
	go gocron.Http_server()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	println(<-ch)
	gocron.Stop(c, wg)
}
