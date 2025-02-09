package main

import (
	"bufio"
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
	exited    chan bool
	Error     error
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
	var err error

	path, args := ParseCommand(j.config.Command)

	go func() {
		for {
			cmd := exec.Command(path, args...)
			if j.config.Umask != nil {
				syscall.Umask(*j.config.Umask)
			}

			cmd.Stdout = os.Stdout
			if j.config.Stdout != nil {
				f, err := os.Create(*j.config.Stdout)
				if err != nil {
					j.Error = err
					return
				}
				cmd.Stdout = f
			}

			if j.config.Stdin != nil {
				f, err := os.Open(*j.config.Stdin)
				if err != nil {
					j.Error = err
					return
				}
				stdin := bufio.NewReader(f)
				cmd.Stdin = stdin
			}

			cmd.Env = j.config.GetEnvString()
			cmd.Dir = j.config.WorkingDir

			j.cmd = cmd
			err = j.cmd.Start()
			if err != nil {
				fmt.Println(err.Error())
				j.Error = err
				return
			}
			j.status = STATUS_RUNNING
			j.isRunning = true
			j.exited = make(chan bool)
			err = j.cmd.Wait()
			close(j.exited)
			fmt.Printf("%d %s: Finished running EXIT: %d\n", os.Getpid(), j.name, j.cmd.ProcessState.ExitCode())
			j.status = STATUS_EXITED

			if !j.isRunning || j.config.Autorestart == AUTORESTART_NEVER {
				break
			}

			if j.isRunning && j.config.Autorestart == AUTORESTART_UNEXPECTED && !j.config.CheckExpectedExitCode(j.cmd.ProcessState.ExitCode()) {
				j.isRunning = false
				return
			}

		}
	}()
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
	signal, exists := SIGNALS[j.config.Stopsignal]
	if exists == false {
		return fmt.Errorf("Unknown signal %s", j.config.Stopsignal)
	}
	j.isRunning = false
	err := j.cmd.Process.Signal(signal)
	if err != nil {
		return err
	}
	select {
	case <-j.exited:
		return nil
	case <-time.After(time.Second * time.Duration(j.config.Stoptime)):
		j.cmd.Process.Kill()
		select {
		case <-j.exited:
			return nil
		case <-time.After(time.Second * 10):
			j.Kill()
			return fmt.Errorf("Process did not exit after timeout. Forced kill.")
		}
	}
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
