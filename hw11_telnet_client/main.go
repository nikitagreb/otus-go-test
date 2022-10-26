package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout")
	flag.Parse()

	host, port := flag.Arg(0), flag.Arg(1)
	if host == "" || port == "" {
		fmt.Fprint(os.Stderr, "server or port is not valid")
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	addr := net.JoinHostPort(host, port)
	cl := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)
	if err := cl.Connect(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	defer cl.Close()

	go func() {
		defer cancel()
		if err := cl.Send(); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			return
		}
	}()

	go func() {
		defer cancel()
		if err := cl.Receive(); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			return
		}
	}()

	<-ctx.Done()
}
