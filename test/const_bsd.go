// +build darwin freebsd

package main

import (
	unix "golang.org/x/sys/unix"
)

var reusePort = unix.SO_REUSEPORT
var reuseAddr = unix.SO_REUSEADDR
