package main

import (
	"net"
	"os"
	"strconv"

	unix "golang.org/x/sys/unix"
)

// Listen returns net.Listener created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func Listen(proto, addr string) (l net.Listener, err error) {
	var (
		soType, fd int
		file       *os.File
		sockaddr   unix.Sockaddr
	)

	if sockaddr, soType, err = getSockaddr(proto, addr); err != nil {
		return nil, err
	}

	if fd, err = unix.Socket(soType, unix.SOCK_STREAM, unix.IPPROTO_TCP); err != nil {
		return nil, err
	}

	if err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, reusePort, 1); err != nil {
		return nil, err
	}

	if err = unix.SetsockoptInt(fd, unix.SOL_SOCKET, reuseAddr, 1); err != nil {
		return nil, err
	}

	if err = unix.Bind(fd, sockaddr); err != nil {
		return nil, err
	}

	// Set backlog size to the maximum
	if err = unix.Listen(fd, unix.SOMAXCONN); err != nil {
		return nil, err
	}

	// File Name get be nil
	file = os.NewFile(uintptr(fd), filePrefix+strconv.Itoa(os.Getpid()))
	if l, err = net.FileListener(file); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return l, err
}
