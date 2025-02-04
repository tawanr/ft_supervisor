package main

import (
	"fmt"
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
	defer s.Close()
	s.Write([]byte(args[0]))

	buf := make([]byte, 1024)
	n, err := s.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[:n]))
}
