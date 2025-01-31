package main

import (
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

type ConfigProgram struct {
	Command     string `yaml:"command"`
	Autorestart bool   `yaml:"autorestart"`
	Autostart   bool   `yaml:"autostart"`
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

func NewConfigParser(filepath string, errorInterface io.Writer) *ConfigParser {
	return &ConfigParser{
		filepath:       filepath,
		errorInterface: errorInterface,
	}
}

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
