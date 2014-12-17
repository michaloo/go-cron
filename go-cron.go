package main

import (
	"bytes"
	"encoding/json"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var build string
var last_err LastRun

type LastRun struct {
	Exit_status int
	Stdout      string
	Stderr      string
	Time        string
	Schedule    string
}

func execute(command string, args []string) {

	log.Println("executing:", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v")
	}

	last_err.Exit_status = 0

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			// so set the error code to tremporary value
			last_err.Exit_status = 127
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				last_err.Exit_status = status.ExitStatus()
				log.Printf("Exit Status: %d", last_err.Exit_status)
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
	last_err.Time = time.Now().Format(time.RFC3339)
	last_err.Stderr = stderr.String()
	last_err.Stdout = stdout.String()
}

func create() (cr *cron.Cron, wgr *sync.WaitGroup) {
	var schedule string = os.Args[1]
	var command string = os.Args[2]
	var args []string = os.Args[3:len(os.Args)]

	last_err.Schedule = schedule

	wg := &sync.WaitGroup{}

	c := cron.New()
	log.Println("new cron:", schedule)

	c.AddFunc(schedule, func() {
		wg.Add(1)
		execute(command, args)
		wg.Done()
	})

	return c, wg
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(last_err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if last_err.Exit_status != 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	w.Write(js)
}

func http_server(c *cron.Cron, wg *sync.WaitGroup) {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":18080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func start(c *cron.Cron, wg *sync.WaitGroup) {
	c.Start()
}

func stop(c *cron.Cron, wg *sync.WaitGroup) {
	log.Println("Stopping")
	c.Stop()
	log.Println("Waiting")
	wg.Wait()
	log.Println("Exiting")
	os.Exit(0)
}

func main() {
	log.Println("Running version: %s", build)

	if len(os.Args) < 3 {
		log.Fatalf("run: go-cron <schedule> <command>")
	}

	c, wg := create()

	go start(c, wg)
	go http_server(c, wg)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	println(<-ch)
	stop(c, wg)
}
