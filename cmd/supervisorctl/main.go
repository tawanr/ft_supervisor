package main

import (
	"net"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		println("usage: supervisorctl [command]")
		return
	}
	s, err := net.Dial("unix", "/tmp/ft_supervisor.sock")
	if err != nil {
		panic(err)
	}
	s.Write([]byte(args[0]))
	s.Close()
}
