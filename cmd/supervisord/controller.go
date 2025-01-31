package main

type Controller struct {
	jobs map[string]*Job
}

func NewController(config Config) *Controller {
	c := &Controller{}
	c.jobs = make(map[string]*Job)
	for name, program := range config.Programs {
		c.jobs[name] = NewJob(&program)
	}
	return c
}

// Start running the controller. This will start all the jobs with enabled autostart
func (c *Controller) Start() error {
	for _, job := range c.jobs {
		if job.config.Autostart == true {
			err := job.Run()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
