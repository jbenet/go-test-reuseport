package main

import (
	"errors"
	"net"
	"syscall"
)

const (
	tcp4                  = 52 // "4"
	tcp6                  = 54 // "6"
	unsupportedProtoError = "Only tcp4 and tcp6 are supported"
	filePrefix            = "port."
)

var reusePort = syscall.SO_REUSEPORT
var reuseAddr = syscall.SO_REUSEADDR

// getSockaddr parses protocol and address and returns implementor syscall.Sockaddr: syscall.SockaddrInet4 or syscall.SockaddrInet6.
func getSockaddr(proto, addr string) (sa syscall.Sockaddr, soType int, err error) {
	var (
		addr4 [4]byte
		addr6 [16]byte
		ip    *net.TCPAddr
	)

	ip, err = net.ResolveTCPAddr(proto, addr)
	if err != nil {
		return nil, -1, err
	}

	switch proto[len(proto)-1] {
	default:
		return nil, -1, errors.New(unsupportedProtoError)
	case tcp4:
		if ip.IP != nil {
			copy(addr4[:], ip.IP[12:16]) // copy last 4 bytes of slice to array
		}
		return &syscall.SockaddrInet4{Port: ip.Port, Addr: addr4}, syscall.AF_INET, nil
	case tcp6:
		if ip.IP != nil {
			copy(addr6[:], ip.IP) // copy all bytes of slice to array
		}
		return &syscall.SockaddrInet6{Port: ip.Port, Addr: addr6}, syscall.AF_INET6, nil
	}
}
