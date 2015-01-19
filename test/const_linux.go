// +build linux

package main

import (
	unix "golang.org/x/sys/unix"
)

var reusePort = 15 // this is not defined in unix go pkg.
var reuseAddr = unix.SO_REUSEADDR
