package main

import (
	"fmt"
	"io"
	"os"

	// ma "github.com/jbenet/go-multiaddr"
	// manet "github.com/jbenet/go-multiaddr/net"
)

func main() {

	// a1 := ma.StringCast("/ip4/0.0.0.0/tcp/11111")
	// a2 := ma.StringCast("/ip4/0.0.0.0/tcp/22222")

	l1, err := Listen("tcp4", "0.0.0.0:11111")
	maybeDie(err)

	l2, err := Listen("tcp4", "0.0.0.0:22222")
	maybeDie(err)

	// // make sure port can be reused. TOOD this doesn't work...
	// err = setSocketReuseListener(l1)
	// maybeDie(err)

	// // make sure port can be reused. TOOD this doesn't work...
	// err = setSocketReuseListener(l2)
	// maybeDie(err)

	// d1 := manet.Dialer{}
	// d2 := manet.Dialer{}
	d1 := dialer{LocalAddr: "127.0.0.1:11111"}
	d2 := dialer{LocalAddr: "127.0.0.1:33333"}

	go func() {
		l2to1foo, err := l2.Accept()
		maybeDie(err)

		fmt.Println("safe")

		l1to2bar, err := l1.Accept()
		maybeDie(err)

		io.Copy(l1to2bar, l2to1foo)
	}()

	d1to2foo, err := d1.Dial("tcp4", "127.0.0.1:22222")
	maybeDie(err)

	d2to1bar, err := d2.Dial("tcp4", "127.0.0.1:11111")
	maybeDie(err)

	go io.Copy(d1to2foo, os.Stdin)
	io.Copy(os.Stdout, d2to1bar)
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(-1)
}

func maybeDie(err error) {
	if err != nil {
		die(err)
	}
}
