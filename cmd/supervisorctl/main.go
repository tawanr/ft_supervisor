package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
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
	s.Write([]byte(strings.Join(args, " ")))

	buf := make([]byte, 1024)
	for {
		n, err := s.Read(buf)
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Print(string(buf[:n]))
	}
}
