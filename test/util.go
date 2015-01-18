package main

import (
	"errors"
	"net"

	unix "golang.org/x/sys/unix"
)

const (
	tcp4                  = 52 // "4"
	tcp6                  = 54 // "6"
	unsupportedProtoError = "Only tcp4 and tcp6 are supported"
	filePrefix            = "port."
)

var reusePort = unix.SO_REUSEPORT
var reuseAddr = unix.SO_REUSEADDR

// getSockaddr parses protocol and address and returns implementor unix.Sockaddr: unix.SockaddrInet4 or unix.SockaddrInet6.
func getSockaddr(proto, addr string) (sa unix.Sockaddr, soType int, err error) {
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
		return &unix.SockaddrInet4{Port: ip.Port, Addr: addr4}, unix.AF_INET, nil
	case tcp6:
		if ip.IP != nil {
			copy(addr6[:], ip.IP) // copy all bytes of slice to array
		}
		return &unix.SockaddrInet6{Port: ip.Port, Addr: addr6}, unix.AF_INET6, nil
	}
}
