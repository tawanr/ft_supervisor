package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Job struct {
	name      string
	config    *ConfigProgram
	pid       int
	cmd       *exec.Cmd
	isRunning bool
	startTime time.Time
	status    JobStatus
}

type JobStatus string

const (
	STATUS_RUNNING  JobStatus = "RUNNING"
	STATUS_STOPPED  JobStatus = "STOPPED"
	STATUS_EXITED   JobStatus = "EXITED"
	STATUS_STARTING JobStatus = "STARTING"
	STATUS_STOPPING JobStatus = "STOPPING"
	STATUS_FATAL    JobStatus = "FATAL"
)

var SIGNALS map[string]os.Signal = map[string]os.Signal{
	"HUP":    syscall.SIGHUP,
	"INT":    syscall.SIGINT,
	"QUIT":   syscall.SIGQUIT,
	"KILL":   syscall.SIGKILL,
	"TERM":   syscall.SIGTERM,
	"USR1":   syscall.SIGUSR1,
	"USR2":   syscall.SIGUSR2,
	"WINCH":  syscall.SIGWINCH,
	"CONT":   syscall.SIGCONT,
	"STOP":   syscall.SIGSTOP,
	"TSTP":   syscall.SIGTSTP,
	"TTIN":   syscall.SIGTTIN,
	"TTOU":   syscall.SIGTTOU,
	"ABRT":   syscall.SIGABRT,
	"ALRM":   syscall.SIGALRM,
	"FPE":    syscall.SIGFPE,
	"ILL":    syscall.SIGILL,
	"TRAP":   syscall.SIGTRAP,
	"BUS":    syscall.SIGBUS,
	"XCPU":   syscall.SIGXCPU,
	"XFSZ":   syscall.SIGXFSZ,
	"SYS":    syscall.SIGSYS,
	"URG":    syscall.SIGURG,
	"IOT":    syscall.SIGIOT,
	"CLD":    syscall.SIGCLD,
	"POLL":   syscall.SIGPOLL,
	"PWR":    syscall.SIGPWR,
	"UNUSED": syscall.SIGUNUSED,
	"IO":     syscall.SIGIO,
	"STKFLT": syscall.SIGSTKFLT,
	"PROF":   syscall.SIGPROF,
}

func NewJob(name string, config *ConfigProgram) *Job {
	return &Job{
		name:      name,
		config:    config,
		isRunning: false,
		status:    STATUS_STOPPED,
	}
}

func (j *Job) Start() error {
	j.pid = j.cmd.Process.Pid
	return nil
}

// Run the job based on configuration
func (j *Job) Run() error {
	path, args := ParseCommand(j.config.Command)
	fmt.Println(path, args)
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	j.cmd = cmd
	err := j.cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	j.status = STATUS_RUNNING
	j.isRunning = true
	err = j.cmd.Wait()
	fmt.Printf("%d %s: Finished running EXIT: %d\n", os.Getpid(), j.name, j.cmd.ProcessState.ExitCode())
	j.status = STATUS_EXITED
	j.isRunning = false
	return err
}

func (j *Job) Status() string {
	text := fmt.Sprintf("%s\t\t%s", j.name, j.status)
	if j.isRunning {
		text += fmt.Sprintf("\t\tpid %d", j.cmd.Process.Pid)
	}
	return text
}

func (j *Job) Exit() error {
	return nil
}

func (j *Job) Kill() error {
	return j.cmd.Process.Kill()
}

// ParseCommand parses a command string into a path and arguments
// Returns command path and arguments
func ParseCommand(command string) (string, []string) {
	quoted := false
	args := strings.FieldsFunc(command, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})
	path, err := exec.LookPath(args[0])
	if err != nil {
		path = args[0]
	}
	return path, args[1:]
}
