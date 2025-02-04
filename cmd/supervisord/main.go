package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	DEFAULT_CONFIG_FILE = "ft_supervisor.yaml"
)

func main() {
	// msg, err := net.Listen("unix", "/tmp/ft_supervisor.sock")
	// defer msg.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// for {
	// 	conn, err := msg.Accept()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	go func(conn net.Conn) {
	// 		txt := make([]byte, 1024)
	// 		l, err := conn.Read(txt)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fmt.Println(string(txt[:l]))
	// 	}(conn)
	// }
	configFile := flag.String("c", DEFAULT_CONFIG_FILE, "Location of the configuration file")
	config := NewConfigParser(*configFile, os.Stderr)
	err := config.Parse()
	if err != nil {
		panic(err.Error())
	}
	err = config.Validate()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(config)
	controller := NewController(config.Config)
	err = controller.Startup()
	if err != nil {
		panic(err.Error())
	}
	for {
		fmt.Print("> ")
		command := ""
		fmt.Scanln(&command)
		exit, err := controllerCommand(controller, command)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if exit == true {
			break
		}
	}
}

func controllerCommand(controller *Controller, command string) (bool, error) {
	switch command {
	case "start":
		return false, controller.Startup()
	case "stop":
		return false, controller.Stop()
	case "status":
		fmt.Println(controller.Status())
		return false, nil
	case "exit":
		controller.Stop()
		return true, nil
	default:
		return false, fmt.Errorf("Unknown command %s", command)
	}
}
