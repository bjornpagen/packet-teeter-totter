package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/davecgh/go-spew/spew"
)

const (
	BufSize = 4096
	PortA   = 12345
	PortB   = 12346
	PortIn  = 12347
)

func isWhitelisted(ipAddr uint32) bool {
	whitelistedIP := binary.BigEndian.Uint32(net.ParseIP("192.168.1.100").To4())
	return ipAddr == whitelistedIP
}

func main() {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create raw TCP socket: %v\n", err)
		os.Exit(1)
	}
	defer syscall.Close(sock)

	destAddr := &syscall.SockaddrInet4{Addr: [4]byte{127, 0, 0, 1}}
	buffer := make([]byte, BufSize)

	for {
		n, srcAddr, err := syscall.Recvfrom(sock, buffer, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to receive packet: %v\n", err)
			continue
		}

		bufBytes := buffer[:n]
		spew.Dump(bufBytes)

		ipHeader := buffer[:20]
		ihl := int(ipHeader[0]&0x0F) << 2
		transportHeader := buffer[ihl : ihl+8]
		dstPort := binary.BigEndian.Uint16(transportHeader[2:4])

		if dstPort != PortIn {
			continue
		}

		srcAddrInet, _ := srcAddr.(*syscall.SockaddrInet4)
		srcIP := binary.BigEndian.Uint32(srcAddrInet.Addr[:])
		if isWhitelisted(srcIP) {
			destAddr.Port = PortA
		} else {
			destAddr.Port = PortB
		}

		err = syscall.Sendto(rawSocket, bufBytes, 0, destAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send packet: %v\n", err)
			continue
		}
	}
}
