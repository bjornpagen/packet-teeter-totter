package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	s := Server{}
	if err := s.Run(); err != nil {
		panic(err)
	}
}

type Server struct {
	buf []byte
}

func (s *Server) Run() error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		return fmt.Errorf("failed to create socket: %w", err)
	}

	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd.%d", fd))

	s.buf = make([]byte, 4096)
	for {
		numRead, err := f.Read(s.buf)
		if err != nil {
			fmt.Println(err)
		}
		bytes := s.buf[:numRead]
		fmt.Printf("% X\n", bytes)
	}
}
