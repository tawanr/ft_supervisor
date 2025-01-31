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
	err = controller.Start()
	if err != nil {
		panic(err.Error())
	}
}
