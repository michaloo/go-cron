package gocron

import (
	"bytes"
	//"fmt"
	"github.com/robfig/cron"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type LastRun struct {
	Exit_status  int
	Stdout       string
	Stderr       string
	ExitTime     string
	Pid          int
	StartingTime string
}

type CurrentState struct {
	Running  map[string]LastRun
	Last     LastRun
	Schedule string
}

var running_processes = map[string]LastRun{}
var Current_state CurrentState

func execute(command string, args []string) {

	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	run := LastRun{}
	run.StartingTime = time.Now().Format(time.RFC3339)

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v")
	}
	run.Pid = cmd.Process.Pid

	log.Println(run.Pid, "cmd:", command, strings.Join(args, " "))

	Current_state.Running[strconv.Itoa(run.Pid)] = run

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			// so set the error code to tremporary value
			run.Exit_status = 127
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				run.Exit_status = status.ExitStatus()
				log.Printf("%d Exit Status: %d", run.Pid, run.Exit_status)
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	log.Printf("%d stdout: %v", run.Pid, stdout.String())
	log.Printf("%d stderr: %v", run.Pid, stderr.String())

	run.ExitTime = time.Now().Format(time.RFC3339)
	run.Stderr = stderr.String()
	run.Stdout = stdout.String()

	delete(Current_state.Running, strconv.Itoa(run.Pid))
	run.Pid = 0
	Current_state.Last = run
}

func Create() (cr *cron.Cron, wgr *sync.WaitGroup) {
	var schedule string = os.Args[1]
	var command string = os.Args[2]
	var args []string = os.Args[3:len(os.Args)]

	wg := &sync.WaitGroup{}

	c := cron.New()
	Current_state = CurrentState{map[string]LastRun{}, LastRun{}, schedule}
	log.Println("new cron:", schedule)

	c.AddFunc(schedule, func() {
		wg.Add(1)
		execute(command, args)
		wg.Done()
	})

	return c, wg
}

func Start(c *cron.Cron) {
	c.Start()
}

func Stop(c *cron.Cron, wg *sync.WaitGroup) {
	log.Println("Stopping")
	c.Stop()
	log.Println("Waiting")
	wg.Wait()
	log.Println("Exiting")
	os.Exit(0)
}
