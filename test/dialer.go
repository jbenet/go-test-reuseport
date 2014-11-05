package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
)

type dialer struct {
	LocalAddr string
}

func (d *dialer) Dial(proto, addr string) (c net.Conn, err error) {
	var (
		soType, fd     int
		file           *os.File
		remoteSockaddr syscall.Sockaddr
		localSockaddr  syscall.Sockaddr
	)

	if remoteSockaddr, soType, err = getSockaddr(proto, addr); err != nil {
		return nil, err
	}

	if localSockaddr, _, err = getSockaddr(proto, d.LocalAddr); err != nil {
		return nil, err
	}

	if fd, err = syscall.Socket(soType, syscall.SOCK_STREAM, syscall.IPPROTO_TCP); err != nil {
		fmt.Println("tcp socket failed")
		return nil, err
	}

	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, reuseAddr, 1); err != nil {
		fmt.Println("reuse addr failed")
		return nil, err
	}

	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, reusePort, 1); err != nil {
		fmt.Println("reuse port failed")
		return nil, err
	}

	if err = syscall.Bind(fd, localSockaddr); err != nil {
		fmt.Println("bind failed")
		return nil, err
	}

	// Set backlog size to the maximum
	if err = syscall.Connect(fd, remoteSockaddr); err != nil {
		fmt.Println("connect failed")
		return nil, err
	}

	// File Name get be nil
	file = os.NewFile(uintptr(fd), filePrefix+strconv.Itoa(os.Getpid()))
	if c, err = net.FileConn(file); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return c, err
}
