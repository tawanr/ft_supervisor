package main

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

func (c *Controller) Stop() error {
	for _, job := range c.jobs {
		err := job.Kill()
		if err != nil {
			return err
		}
	}
	return nil
}
