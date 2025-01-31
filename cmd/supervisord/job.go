package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"text/scanner"
)

type Job struct {
	config    *ConfigProgram
	pid       int
	cmd       *exec.Cmd
	isRunning bool
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

func NewJob(config *ConfigProgram) *Job {
	return &Job{
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
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	j.cmd = cmd
	err := j.cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = j.cmd.Wait()
	return err
}

// ParseCommand parses a command string into a path and arguments
// Returns command path and arguments
func ParseCommand(command string) (string, []string) {
	var s scanner.Scanner
	s.Init(strings.NewReader(command))
	args := []string{}
	token := s.Scan()
	for token != scanner.EOF {
		args = append(args, s.TokenText())
		token = s.Scan()
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		path = args[0]
	}
	return path, args[1:]
}
