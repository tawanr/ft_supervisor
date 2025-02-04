package main

import (
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

type ConfigProgram struct {
	Command     string  `yaml:"command" validate:"required"`
	Autorestart bool    `yaml:"autorestart"`
	Autostart   bool    `yaml:"autostart"`
	Exitcodes   []int   `yaml:"exitcodes"`
	Stopsignal  string  `yaml:"stopsignal"`
	Stoptime    int     `yaml:"stoptime"`
	Numprocs    int     `yaml:"numprocs"`
	Stdin       *string `yaml:"stdin"`
	Stdout      *string `yaml:"stdout"`
}

type Config struct {
	Programs map[string]ConfigProgram `yaml:"programs"`
}

type ConfigParser struct {
	Config         Config
	filepath       string
	errors         []error
	errorInterface io.Writer
}

func (p *ConfigProgram) setDefaults() {
	p.Autorestart = true
	p.Autostart = false
	p.Exitcodes = []int{0}
	p.Stopsignal = "TERM"
	p.Stoptime = 0
	p.Numprocs = 1
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (p *ConfigProgram) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Declaring a new type so that the custom unmarshaller for the current type is not called infinitely.
	type plain ConfigProgram
	p.setDefaults()
	return unmarshal((*plain)(p))
}

func NewConfigParser(filepath string, errorInterface io.Writer) *ConfigParser {
	return &ConfigParser{
		filepath:       filepath,
		errorInterface: errorInterface,
	}
}

// Parse parses the YAML config file and returns an error if there is one.
func (c *ConfigParser) Parse() error {
	file, err := os.Open(c.filepath)
	contents, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(contents, &c.Config)
	if err != nil {
		return err
	}
	err = c.Validate()
	if err != nil {
		c.printError()
		return err
	}
	return nil
}

// Validate validates the config file and returns an error if there is one.
func (c *ConfigParser) Validate() error {
	for name, program := range c.Config.Programs {
		if program.Command == "" {
			c.errors = append(c.errors, fmt.Errorf("Program %s has no command", name))
		}
	}
	if len(c.errors) > 0 {
		return fmt.Errorf("Config file has errors")
	}
	return nil
}

func (c *ConfigParser) printError() {
	for _, err := range c.errors {
		fmt.Fprintf(c.errorInterface, "Error: %s\n", err.Error())
	}
}
