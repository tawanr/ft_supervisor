package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

const (
	DEFAULT_CONFIG_FILE = "ft_supervisor.conf"
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
	file, err := os.Open(*configFile)
	if err != nil {
		panic(err)
	}
	contents, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	config := Config{}
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		panic(err)
	}
	err = config.Validate()
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(config)
}
