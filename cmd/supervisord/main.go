package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

const (
	DEFAULT_CONFIG_FILE = "ft_supervisor.yaml"
)

func main() {
	configFile := flag.String("c", DEFAULT_CONFIG_FILE, "Location of the configuration file")
	interactive := flag.Bool("i", false, "Interactive mode")
	flag.Parse()
	config := NewConfigParser(*configFile, os.Stderr)
	err := config.Parse()
	if err != nil {
		panic(err.Error())
	}
	err = config.Validate()
	if err != nil {
		panic(err.Error())
	}
	controller := NewController(config.Config)
	err = controller.Startup()
	if err != nil {
		panic(err.Error())
	}
	if *interactive == true {
		promptInput(controller)
	} else {
		socketInput(controller)
	}
}

func socketInput(controller *Controller) {
	cmd, err := net.Listen("unix", "/tmp/ft_supervisor.sock")
	if err != nil {
		panic(err)
	}
	defer cmd.Close()
	for {
		conn, err := cmd.Accept()
		if err != nil {
			panic(err)
		}

		txt := make([]byte, 1024)
		l, err := conn.Read(txt)
		if err != nil {
			panic(err)
		}
		output := ""
		exit, err := controllerCommand(controller, string(txt[:l]), &output)
		if err != nil {
			panic(err)
		}
		if len(output) > 0 {
			conn.Write([]byte(output))
		}
		if exit == true {
			return
		}
		conn.Close()
	}
}

func promptInput(controller *Controller) {
	for {
		fmt.Print("> ")
		command := ""
		fmt.Scanln(&command)
		output := ""
		exit, err := controllerCommand(controller, command, &output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if len(output) > 0 {
			fmt.Println(output)
		}
		if exit == true {
			return
		}
	}
}

func controllerCommand(controller *Controller, command string, output *string) (bool, error) {
	switch command {
	case "start":
		return false, controller.Startup()
	case "stop":
		return false, controller.Stop()
	case "status":
		*output = controller.Status()
		return false, nil
	case "exit":
		controller.Stop()
		return true, nil
	default:
		return false, fmt.Errorf("Unknown command %s", command)
	}
}
