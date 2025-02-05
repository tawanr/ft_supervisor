package main

import "fmt"

type Controller struct {
	jobs map[string]*Job
}

func NewController(config Config) *Controller {
	c := &Controller{}
	c.jobs = make(map[string]*Job)
	for name, program := range config.Programs {
		c.jobs[name] = NewJob(name, &program)
	}
	return c
}

// Start running the controller. This will start all the jobs with enabled autostart
func (c *Controller) Startup() error {
	for _, job := range c.jobs {
		if job.config.Autostart == true {
			go job.Run()
		}
	}
	return nil
}

func (c *Controller) Status() string {
	status := ""
	for _, job := range c.jobs {
		status += job.Status() + "\n"
	}
	return status
}

func (c *Controller) Start(programs []string) error {
	startingJobs := []*Job{}
	if len(programs) == 1 && programs[0] == "all" {
		for _, job := range c.jobs {
			startingJobs = append(startingJobs, job)
		}
	} else {
		for _, program := range programs {
			_, exists := c.jobs[program]
			if exists == false {
				return fmt.Errorf("Job %s not found", program)
			}
			startingJobs = append(startingJobs, c.jobs[program])
		}
	}
	for _, program := range startingJobs {
		go program.Run()
	}
	return nil
}

func (c *Controller) Stop(programs []string) error {
	stoppingJobs := []*Job{}
	if len(programs) == 1 && programs[0] == "all" {
		for _, job := range c.jobs {
			stoppingJobs = append(stoppingJobs, job)
		}
	} else {
		for _, program := range programs {
			_, exists := c.jobs[program]
			if exists == false {
				return fmt.Errorf("Job %s not found", program)
			}
			stoppingJobs = append(stoppingJobs, c.jobs[program])
		}
	}
	for _, program := range stoppingJobs {
		// TODO: Error handling
		go program.Exit()
	}
	return nil
}

func (c *Controller) Exit() error {
	for _, job := range c.jobs {
		err := job.Exit()
		if err != nil {
			return err
		}
	}
	return nil
}
